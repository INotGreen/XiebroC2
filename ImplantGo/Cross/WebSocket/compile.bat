go env -w GOPROXY=https://goproxy.cn,direct
go get github.com/andreburgaud/crypt2go@v1.4.1
set GOOS=linux&&go build -a -ldflags="-s -w" -installsuffix cgo -o wslMain main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o wsMac_AMD_x64 main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o wsMac_ARM_x64 main.go
