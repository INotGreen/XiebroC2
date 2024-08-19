package TCPsocket

import (
	"encoding/binary"
	"main/Encrypt"
	"net"
)

func Send(conn net.Conn, msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			//log.Println("Send error:", err)
		}
	}()

	if conn == nil {
		//log.Println("Connection not established")
		return
	}
	msg, err := Encrypt.Encrypt(msg)
	if err != nil {
		return
	}
	bufferSize := int32(len(msg))
	err = binary.Write(conn, binary.LittleEndian, bufferSize)
	if err != nil {
		//log.Println("Failed to send buffer size:", err)
		return
	}

	const chunkSize = 50 * 1024 // 50 KB
	var chunk []byte

	for bytesSent := 0; bytesSent < len(msg); {
		if len(msg)-bytesSent > chunkSize {
			chunk = msg[bytesSent : bytesSent+chunkSize]
		} else {
			chunk = msg[bytesSent:]
		}

		_, err := conn.Write(chunk)
		if err != nil {
			//log.Println("Failed to send data:", err)
			return
		}

		bytesSent += len(chunk)
	}
}
