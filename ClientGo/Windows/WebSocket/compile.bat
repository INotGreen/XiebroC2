
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wsMain.exe main.go