set GOOS=windows
set GOARCH=amd64
go build -a -ldflags="-s -w" -installsuffix cgo main.go&&main.exe