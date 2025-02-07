
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wsMain_amd64.exe cmd\ws\windows\main.go
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o tcpMain_amd64.exe cmd\tcp\windows\main.go
set GOOS=linux&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wslMain_amd64 cmd\ws\linux\main.go
set GOOS=linux&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o tcplMain_amd64 cmd\tcp\linux\main.go