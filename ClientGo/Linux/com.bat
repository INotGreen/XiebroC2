set GOOS=linux
go build -a -ldflags="-s -w" -installsuffix cgo main.go