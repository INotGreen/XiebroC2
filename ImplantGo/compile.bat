
amd64
set GOOS=windows&&set GOARCH=amd64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Main_amd64.exe cmd\windows\main.go
set GOOS=linux&&set GOARCH=amd64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Linux_amd64 cmd\linux\main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Darwin_amd64 cmd\darwin\main.go




set GOOS=windows&&set GOARCH=arm64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Main_arm64.exe cmd\windows\main.go
set GOOS=linux&&set GOARCH=arm64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Linux_arm64 cmd\linux\main.go
set GOOS=darwin&&set GOARCH=arm64&&go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Darwin_arm64 cmd\darwin\main.go


