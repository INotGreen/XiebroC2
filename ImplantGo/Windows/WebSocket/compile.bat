
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wsMain_amd64.exe cmd\Windows\main.go
set GOOS=linux&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wsMain_amd64. cmd\Linux\main.go