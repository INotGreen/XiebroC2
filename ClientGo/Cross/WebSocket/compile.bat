
go mod init main
go env -w GOPROXY=https://goproxy.cn,direct
go get golang.org/x/text/encoding/unicode
go get golang.org/x/text/encoding/unicode
go get github.com/Binject/debug/pe
go get github.com/Ne0nd0g/go-clr
go get github.com/andreburgaud/crypt2go@v1.4.1
go get github.com/Binject/go-donut/donut
go get github.com/kbinani/screenshot
go get github.com/shirou/gopsutil/net
go get github.com/shirou/gopsutil/process
go get github.com/togettoyou/wsc
go get github.com/xtaci/smux
go get github.com/creack/pty
go get golang.org/x/term
set GOOS=linux&&go build -a -ldflags="-s -w" -installsuffix cgo -o wslMain main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o Intel_x86_64 main.go
set GOOS=darwin&&set GOARCH=amd64&&go build -a -ldflags="-s -w" -installsuffix cgo -o M1_x86_64 main.go
