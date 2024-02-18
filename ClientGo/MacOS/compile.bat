set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o Intel_x86_64 main.go

set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o M1_x86_64 main.go
