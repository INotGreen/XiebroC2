package mux

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"golang.org/x/crypto/chacha20poly1305"
)

type frameHeader struct {
	id     uint32
	length uint16
	flags  uint16
	nonce  [chacha20poly1305.NonceSizeX]byte
}

const (
	frameHeaderSize = 4 + 2 + 2 + chacha20poly1305.NonceSizeX // must be exact frameHeader struct size
	maxPayloadSize  = math.MaxUint16
)

const (
	flagData        = iota + 1 // data frame
	flagKeepalive              // empty frame to keep connection open
	flagOpenStream             // first frame in stream
	flagCloseStream            // stream is being closed gracefully
	flagCloseMux               // mux is being closed gracefully
)

func encodeFrameHeader(buf []byte, h frameHeader) {
	binary.LittleEndian.PutUint32(buf[0:], h.id)
	binary.LittleEndian.PutUint16(buf[4:], h.length)
	binary.LittleEndian.PutUint16(buf[6:], h.flags)
	copy(buf[8:], h.nonce[:])
}

func decodeFrameHeader(buf []byte) frameHeader {
	h := frameHeader{
		id:     binary.LittleEndian.Uint32(buf[0:]),
		length: binary.LittleEndian.Uint16(buf[4:]),
		flags:  binary.LittleEndian.Uint16(buf[6:]),
	}
	copy(h.nonce[:], buf[8:])
	return h
}

// nextFrame reads and decrypts a frame from reader
func readFrame(reader io.Reader, aead cipher.AEAD, frameBuf []byte) (frameHeader, []byte, error) {

	headerBuf := [frameHeaderSize]byte{}

	if _, err := io.ReadFull(reader, headerBuf[:]); err != nil {
		return frameHeader{}, nil, fmt.Errorf("could not read frame header: %w", err)
	}
	h := decodeFrameHeader(headerBuf[:])

	payloadSize := uint32(chacha20poly1305.Overhead)
	if h.flags == flagData {
		payloadSize += uint32(h.length)
	}

	if _, err := io.ReadFull(reader, frameBuf[:payloadSize]); err != nil {
		return frameHeader{}, nil, fmt.Errorf("could not read frame payload: %w", err)
	}

	// Decrypt the message and check it wasn't tampered with.
	if _, err := aead.Open(frameBuf[:0], h.nonce[:], frameBuf[:payloadSize], headerBuf[:]); err != nil {
		return frameHeader{}, nil, err
	}

	return h, frameBuf[:h.length], nil
}

// appendFrame writs and encrypts a frame to buf
func appendFrame(buf []byte, aead cipher.AEAD, h frameHeader, payload []byte) []byte {
	rand.Read(h.nonce[:])
	frame := buf[len(buf):][:frameHeaderSize+len(payload)+aead.Overhead()]
	encodeFrameHeader(frame[:frameHeaderSize], h)
	aead.Seal(frame[frameHeaderSize:][:0], h.nonce[:], payload, frame[:frameHeaderSize])
	return buf[:len(buf)+len(frame)]
}
