go get github.com/andreburgaud/crypt2go@v1.4.1
set GOOS=linux&&go build -a -ldflags="-s -w" -installsuffix cgo -o LinuxMain main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o Mac_AMD_x64 main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o Mac_ARM_x64 main.go