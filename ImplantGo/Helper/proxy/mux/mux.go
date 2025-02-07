package mux

import (
	"crypto/cipher"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
)

// NOTE: This package makes heavy use of sync.Cond to manage concurrent streams
// multiplexed onto a single connection. sync.Cond is rarely used (since it is
// almost never the right tool for the job), and consequently, Go programmers
// tend to be unfamiliar with its semantics. Nevertheless, it is currently the
// only way to achieve optimal throughput in a stream multiplexer, so we make
// careful use of it here. Please make sure you understand sync.Cond thoroughly
// before attempting to modify this code.

// Errors relating to stream or mux shutdown.
var (
	ErrClosedConn       = errors.New("underlying connection was closed")
	ErrClosedStream     = errors.New("stream was gracefully closed")
	ErrPeerClosedStream = errors.New("peer closed stream gracefully")
	ErrPeerClosedConn   = errors.New("peer closed mux gracefully")
	ErrWriteClosed      = errors.New("write end of stream closed")
)

// A Mux multiplexes multiple duplex Streams onto a single net.Conn.
type Mux struct {
	conn       net.Conn
	aead       cipher.AEAD
	acceptChan chan *Stream

	readMutex sync.Mutex
	// subsequent fields are used by readLoop() and guarded by readMutex
	readErr error
	streams map[uint32]*Stream
	nextID  uint32

	writeMutex sync.Mutex
	// subsequent fields are used by writeLoop() and guarded by writeMutex
	writeErr   error
	writeCond  sync.Cond
	bufferCond sync.Cond
	writeBuf   []byte
	sendBuf    []byte
	writeBufA  []byte
	writeBufB  []byte
}

// setErr sets the Mux error and wakes up all Mux-related goroutines. If m.err
// is already set, setErr is a no-op.
func (m *Mux) setErr(err error) error {

	// Set readErr
	m.readMutex.Lock()
	defer m.readMutex.Unlock()
	if m.readErr != nil {
		return m.readErr
	}
	m.readErr = err

	// Set writeErr
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	if m.writeErr != nil {
		return m.writeErr
	}
	m.writeErr = err

	for _, s := range m.streams {
		s.cond.L.Lock()
		s.err = err
		s.cond.L.Unlock()
		s.cond.Broadcast()
	}
	m.conn.Close()
	m.writeCond.Broadcast()
	m.bufferCond.Broadcast()
	close(m.acceptChan)
	return err
}

// bufferFrame blocks until it can store its frame in m.writeBuf. It returns
// early with an error if m.err is set.
func (m *Mux) bufferFrame(h frameHeader, payload []byte) error {
	m.writeMutex.Lock()

	// block until we can add the frame to the buffer
	for len(m.writeBuf)+frameHeaderSize+len(payload)+chacha20poly1305.Overhead > cap(m.writeBuf) && m.writeErr == nil {
		m.bufferCond.Wait()
	}
	if m.writeErr != nil {
		m.writeMutex.Unlock()
		return m.writeErr
	}

	// queue our frame
	m.writeBuf = appendFrame(m.writeBuf, m.aead, h, payload)
	m.writeMutex.Unlock()

	// wake the writeLoop
	m.writeCond.Signal()
	// wake at most one bufferFrame call
	m.bufferCond.Signal()
	return nil
}

// writeLoop handles the actual Writes to the Mux's net.Conn. It waits for
// bufferFrame calls to fill m.writeBuf, then flushes the buffer to the
// underlying connection. It also handles keepalives.
func (m *Mux) writeLoop() {

	keepaliveInterval := time.Minute * 4
	nextKeepalive := time.Now().Add(keepaliveInterval)
	timer := time.AfterFunc(keepaliveInterval, m.writeCond.Signal)
	defer timer.Stop()

	for {

		m.writeMutex.Lock()
		for len(m.writeBuf) == 0 && m.writeErr == nil && time.Now().Before(nextKeepalive) {
			m.writeCond.Wait()
		}

		if m.writeErr != nil {
			m.writeMutex.Unlock()
			return
		}

		// if we have a normal frame, use that; otherwise, send a keepalive
		//
		// NOTE: even if we were woken by the keepalive timer, there might be a
		// normal frame ready to send, in which case we don't need a keepalive
		if len(m.writeBuf) == 0 {
			m.writeBuf = appendFrame(m.writeBuf[:0], m.aead, frameHeader{flags: flagKeepalive}, nil)
		}

		// to avoid blocking bufferFrame while we Write, swap writeBufA and writeBufB
		m.writeBuf, m.sendBuf = m.sendBuf, m.writeBuf
		m.writeMutex.Unlock()

		// wake at most one bufferFrame call
		m.bufferCond.Signal()

		// reset keepalive timer
		timer.Stop()
		timer.Reset(keepaliveInterval)
		nextKeepalive = time.Now().Add(keepaliveInterval)

		// write the packet(s)
		if _, err := m.conn.Write(m.sendBuf); err != nil {
			m.setErr(err)
			return
		}

		// clear sendBuf
		m.sendBuf = m.sendBuf[:0]
	}
}

// Delete stream from Mux
func (m *Mux) deleteStream(id uint32) {
	m.readMutex.Lock()
	delete(m.streams, id)
	m.readMutex.Unlock()
}

// readLoop handles the actual Reads from the Mux's net.Conn. It waits for a
// frame to arrive, then routes it to the appropriate Stream, creating a new
// Stream if required.
func (m *Mux) readLoop() {

	frameBuf := make([]byte, maxPayloadSize+chacha20poly1305.Overhead)

	for {
		header, payload, err := readFrame(m.conn, m.aead, frameBuf)

		if err != nil {
			m.setErr(err)
			return
		}

		switch header.flags {

		case flagKeepalive:
			continue

		case flagOpenStream:
			m.readMutex.Lock()
			s := newStream(header.id, m)
			m.streams[header.id] = s
			m.readMutex.Unlock()
			m.acceptChan <- s

		case flagCloseMux:
			m.setErr(ErrPeerClosedConn)
			return

		default:
			m.readMutex.Lock()
			stream, found := m.streams[header.id]
			m.readMutex.Unlock()
			if found {
				stream.consumeFrame(header, payload)
			} else {
				log.Printf("can't find stream ID (%v) (length=%v, flags=%v)", header.id, header.length, header.flags)
			}
		}
	}
}

// Close closes the underlying net.Conn.
func (m *Mux) Close() error {
	// tell perr we are shutting down
	h := frameHeader{flags: flagCloseMux}
	m.bufferFrame(h, nil)

	// if there's a buffered Write, wait for it to be sent
	m.writeMutex.Lock()
	for len(m.writeBuf) > 0 && m.writeErr == nil {
		m.bufferCond.Wait()
	}
	m.writeMutex.Unlock()
	err := m.setErr(ErrClosedConn)
	if err == ErrClosedConn {
		return nil
	}
	return err
}

// AcceptStream waits for and returns the next peer-initiated Stream.
func (m *Mux) AcceptStream() (net.Conn, error) {
	if s, ok := <-m.acceptChan; ok {
		return s, nil
	}
	return nil, m.readErr
}

// OpenStream creates a new Stream.
func (m *Mux) OpenStream() (net.Conn, error) {
	m.readMutex.Lock()
	s := newStream(m.nextID, m)
	m.streams[s.id] = s
	m.nextID += 2 // int wraparound intended
	m.readMutex.Unlock()

	// send flagOpenStream to tell peer the stream exists
	h := frameHeader{
		id:    s.id,
		flags: flagOpenStream,
	}
	return s, m.bufferFrame(h, nil)
}

// newMux initializes a Mux and spawns its readLoop and writeLoop goroutines.
func newMux(conn net.Conn, startID uint32, psk string) *Mux {
	m := &Mux{
		conn:       conn,
		acceptChan: make(chan *Stream, 256),
		streams:    make(map[uint32]*Stream),
		nextID:     startID,
		writeBufA:  make([]byte, 0, maxPayloadSize*10),
		writeBufB:  make([]byte, 0, maxPayloadSize*10),
	}
	key := blake2b.Sum256([]byte(psk))
	m.aead, _ = chacha20poly1305.NewX(key[:])
	m.writeCond.L = &m.writeMutex  // both conds use the same mutex
	m.bufferCond.L = &m.writeMutex //
	m.writeBuf = m.writeBufA       // initial writeBuf is writeBufA
	m.sendBuf = m.writeBufB        // initial sendBuf is writeBufB

	go m.readLoop()
	go m.writeLoop()
	return m
}

// Client creates and initializes a new client-side Mux on the provided conn.
// Client takes overship of the conn.
func Client(conn net.Conn, psk string) *Mux {
	return newMux(conn, 0, psk)
}

// Server creates and initializes a new server-side Mux on the provided conn.
// Server takes overship of the conn.
func Server(conn net.Conn, psk string) *Mux {
	return newMux(conn, 1, psk)
}
