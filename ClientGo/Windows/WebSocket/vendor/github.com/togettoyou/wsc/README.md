# golang websocket client
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/togettoyou/wsc)
[![GoDoc](https://godoc.org/github.com/togettoyou/wsc?status.svg)](https://godoc.org/github.com/togettoyou/wsc)

### Install

```
$ go get -v github.com/togettoyou/wsc
```
#### Simple example

``` go
package main

import (
	"github.com/togettoyou/wsc"
	"log"
	"time"
)

func main() {
	done := make(chan bool)
	ws := wsc.New("ws://127.0.0.1:7777/ws")
	// 可自定义配置，不使用默认配置
	//ws.SetConfig(&wsc.Config{
	//	// 写超时
	//	WriteWait: 10 * time.Second,
	//	// 支持接受的消息最大长度，默认512字节
	//	MaxMessageSize: 2048,
	//	// 最小重连时间间隔
	//	MinRecTime: 2 * time.Second,
	//	// 最大重连时间间隔
	//	MaxRecTime: 60 * time.Second,
	//	// 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
	//	RecFactor: 1.5,
	//	// 消息发送缓冲池大小，默认256
	//	MessageBufferSize: 1024,
	//})
	// 设置回调处理
	ws.OnConnected(func() {
		log.Println("OnConnected: ", ws.WebSocket.Url)
		// 连接成功后，测试每5秒发送消息
		go func() {
			t := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-t.C:
					err := ws.SendTextMessage("hello")
					if err == wsc.CloseErr {
						return
					}
				}
			}
		}()
	})
	ws.OnConnectError(func(err error) {
		log.Println("OnConnectError: ", err.Error())
	})
	ws.OnDisconnected(func(err error) {
		log.Println("OnDisconnected: ", err.Error())
	})
	ws.OnClose(func(code int, text string) {
		log.Println("OnClose: ", code, text)
		done <- true
	})
	ws.OnTextMessageSent(func(message string) {
		log.Println("OnTextMessageSent: ", message)
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		log.Println("OnBinaryMessageSent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		log.Println("OnSentError: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		log.Println("OnPingReceived: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		log.Println("OnPongReceived: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		log.Println("OnTextMessageReceived: ", message)
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		log.Println("OnBinaryMessageReceived: ", string(data))
	})
	// 开始连接
	go ws.Connect()
	for {
		select {
		case <-done:
			return
		}
	}
}
```
