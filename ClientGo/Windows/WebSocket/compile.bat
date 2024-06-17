
set GOOS=windows&&go build -a -ldflags="-s -w" -installsuffix cgo -o wsMain.exe main.go