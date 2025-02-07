package mux

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// A Stream is a duplex connection multiplexed over a net.Conn. It implements
// the net.Conn interface.
type Stream struct {
	mux     *Mux
	id      uint32
	cond    sync.Cond // guards + synchronizes subsequent fields
	err     error
	readBuf []byte
	rd, wd  time.Time // deadlines
}

func newStream(id uint32, m *Mux) *Stream {
	return &Stream{
		mux:  m,
		id:   id,
		cond: sync.Cond{L: new(sync.Mutex)},
		err:  m.readErr,
	}
}

// LocalAddr returns the underlying connection's LocalAddr.
func (s *Stream) LocalAddr() net.Addr { return s.mux.conn.LocalAddr() }

// RemoteAddr returns the underlying connection's RemoteAddr.
func (s *Stream) RemoteAddr() net.Addr { return s.mux.conn.RemoteAddr() }

// SetDeadline sets the read and write deadlines associated with the Stream. It
// is equivalent to calling both SetReadDeadline and SetWriteDeadline.
//
// This implementation does not entirely conform to the net.Conn interface:
// setting a new deadline does not affect pending Read or Write calls, only
// future calls.
func (s *Stream) SetDeadline(t time.Time) error {
	s.SetReadDeadline(t)
	s.SetWriteDeadline(t)
	return nil
}

// SetReadDeadline sets the read deadline associated with the Stream.
//
// This implementation does not entirely conform to the net.Conn interface:
// setting a new deadline does not affect pending Read calls, only future calls.
func (s *Stream) SetReadDeadline(t time.Time) error {
	s.cond.L.Lock()
	s.rd = t
	s.cond.L.Unlock()
	return nil
}

// SetWriteDeadline sets the write deadline associated with the Stream.
//
// This implementation does not entirely conform to the net.Conn interface:
// setting a new deadline does not affect pending Write calls, only future
// calls.
func (s *Stream) SetWriteDeadline(t time.Time) error {
	s.cond.L.Lock()
	s.wd = t
	s.cond.L.Unlock()
	return nil
}

// consumeFrame processes a frame based on h.flags.
func (s *Stream) consumeFrame(h frameHeader, payload []byte) {
	s.cond.L.Lock()

	switch h.flags {

	case flagCloseStream:
		s.err = ErrPeerClosedStream
		s.mux.deleteStream(s.id)

	case flagData:
		s.readBuf = payload
		s.cond.Broadcast() // wake Read
		for len(s.readBuf) > 0 && s.err == nil && (s.rd.IsZero() || time.Now().Before(s.rd)) {
			s.cond.Wait()
		}

	default:
		// The flags are mutually exclusive, we should never be here
		// ignore as the peer sent a bad frame
		log.Printf("peer sent invalid frame ID (%v) (length=%v, flags=%v)", h.id, h.length, h.flags)

	}
	s.cond.L.Unlock()
}

// Read reads data from the Stream.
func (s *Stream) Read(p []byte) (int, error) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	if !s.rd.IsZero() {
		if !time.Now().Before(s.rd) {
			return 0, os.ErrDeadlineExceeded
		}
		timer := time.AfterFunc(time.Until(s.rd), s.cond.Broadcast)
		defer timer.Stop()
	}

	// Wait for data, an error, stream close, or timeout.
	for len(s.readBuf) == 0 && s.err == nil && (s.rd.IsZero() || time.Now().Before(s.rd)) {
		s.cond.Wait()
	}

	// A sender could have sent some data then closed the stream. We want to
	// return the data before indicating the closure of the stream to keep the
	// order of events correct.
	if len(s.readBuf) > 0 {
		n := copy(p, s.readBuf)
		s.readBuf = s.readBuf[n:]
		s.cond.Broadcast() // wake consumeFrame
		return n, nil
	}

	// Check for errors. There is a very unlikely chance that between leaving
	// the for loop and checking the deadline it expires. We skip checking the
	// deadline if data has been received.

	if s.err != nil && s.err == ErrPeerClosedStream {
		return 0, io.EOF
	} else if s.err != nil {
		return 0, s.err
	} else if !s.rd.IsZero() && !time.Now().Before(s.rd) {
		return 0, os.ErrDeadlineExceeded
	}
	return 0, nil

}

// Write writes data to the Stream.
func (s *Stream) Write(p []byte) (int, error) {

	buf := bytes.NewBuffer(p)
	for buf.Len() > 0 {

		if s.err != nil {
			return len(p) - buf.Len(), s.err
		}

		payload := buf.Next(maxPayloadSize)
		h := frameHeader{
			id:     s.id,
			length: uint16(len(payload)),
			flags:  flagData,
		}

		// write next frame's worth of data
		if err := s.mux.bufferFrame(h, payload); err != nil {
			return len(p) - buf.Len() - len(payload), err
		}
	}
	return len(p), nil
}

// Close closes the Stream. The underlying connection is not closed.
func (s *Stream) Close() error {
	// cancel outstanding Read/Write calls
	//
	// NOTE: Read calls will be interrupted immediately, but Write calls might
	// send another frame before observing the Close. This is ok: the peer will
	// discard any frames that arrive after the flagLast frame.
	s.cond.L.Lock()
	if s.err == ErrClosedStream || s.err == ErrPeerClosedStream {
		s.cond.L.Unlock()
		return nil
	}

	s.err = ErrClosedStream
	s.cond.L.Unlock()
	s.cond.Broadcast()

	h := frameHeader{
		id:    s.id,
		flags: flagCloseStream,
	}
	err := s.mux.bufferFrame(h, nil)
	if err != nil && err != ErrPeerClosedStream {
		return err
	}

	// delete stream from Mux
	s.mux.deleteStream(s.id)
	return nil
}

var _ net.Conn = (*Stream)(nil)
