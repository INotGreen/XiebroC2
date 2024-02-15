package main

import (
	"io"
	"net"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

func GetInteractiveShell(conn net.Conn) {
	paramBuffer := make([]byte, 128)

	paramLen, err := conn.Read(paramBuffer)

	if err != nil {
		return
	}

	termEnv := strings.TrimSpace(string(paramBuffer[:paramLen]))

	_, err = conn.Read(paramBuffer)

	if err != nil {
		return
	}

	var ws pty.Winsize

	ws.Rows = uint16(paramBuffer[0]) + uint16(paramBuffer[1]<<8)
	ws.Cols = uint16(paramBuffer[2]) + uint16(paramBuffer[3]<<8)

	ws.X = 0
	ws.Y = 0

	c := exec.Command("/bin/sh", "-c", "exec bash --login")

	c.Env = append(c.Env, "HISTFILE=/dev/null")

	c.Env = append(c.Env, "TERM="+termEnv)

	ptmx, err := pty.Start(c)
	if err != nil {
		return
	}
	defer func() { _ = ptmx.Close() }()

	_ = pty.Setsize(ptmx, &ws)

	go func() {
		for {
			buff := make([]byte, 1024)
			readLen, err := conn.Read(buff)
			if err != nil {
				break
			}
			if readLen > 0 {
				_, err = ptmx.Write(buff[:readLen])
				if err != nil {
					break
				}
			}
		}
		//不需要对输入做控制的可以直接采用下面的方式
		//_, _ = io.Copy(ptmx, conn)
	}()

	_, _ = io.Copy(conn, ptmx)

}
