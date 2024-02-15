package generate

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/Binject/go-donut/donut"
)

// DonutFromAssembly - Generate a donut shellcode from a .NET assembly
func DonutFromAssembly(assembly []byte, isDLL bool, arch string, params string, method string, className string, appDomain string) ([]byte, error) {
	ext := ".exe"
	if isDLL {
		ext = ".dll"
	}
	donutArch := getDonutArch(arch)
	config := donut.DefaultConfig()
	config.Bypass = 3
	config.Runtime = "v4.0.30319" // we might want to make this configurable
	config.Format = 1
	config.Arch = donutArch
	config.Class = className
	config.Parameters = params
	config.Domain = appDomain
	config.Method = method
	config.Entropy = 3
	config.Unicode = 0
	config.Type = getDonutType(ext, true)
	return getDonut(assembly, config)
}

func getDonut(data []byte, config *donut.DonutConfig) (shellcode []byte, err error) {
	buf := bytes.NewBuffer(data)
	res, err := donut.ShellcodeFromBytes(buf, config)
	if err != nil {
		return
	}
	shellcode = res.Bytes()
	stackCheckPrologue := []byte{
		// Check stack is 8 byte but not 16 byte aligned or else errors in LoadLibrary
		0x48, 0x83, 0xE4, 0xF0, // and rsp,0xfffffffffffffff0
		0x48, 0x83, 0xC4, 0x08, // add rsp,0x8
	}
	shellcode = append(stackCheckPrologue, shellcode...)
	return
}

func getDonutArch(arch string) donut.DonutArch {
	var donutArch donut.DonutArch
	switch strings.ToLower(arch) {
	case "x32", "386":
		donutArch = donut.X32
	case "x64", "amd64":
		donutArch = donut.X64
	case "x84":
		donutArch = donut.X84
	default:
		donutArch = donut.X84
	}
	return donutArch
}

func getDonutType(ext string, dotnet bool) donut.ModuleType {
	var donutType donut.ModuleType
	switch strings.ToLower(filepath.Ext(ext)) {
	case ".exe", ".bin":
		if dotnet {
			donutType = donut.DONUT_MODULE_NET_EXE
		} else {
			donutType = donut.DONUT_MODULE_EXE
		}
	case ".dll":
		if dotnet {
			donutType = donut.DONUT_MODULE_NET_DLL
		} else {
			donutType = donut.DONUT_MODULE_DLL
		}
	case ".xsl":
		donutType = donut.DONUT_MODULE_XSL
	case ".js":
		donutType = donut.DONUT_MODULE_JS
	case ".vbs":
		donutType = donut.DONUT_MODULE_VBS
	}
	return donutType
}
