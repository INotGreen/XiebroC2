package args

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"main/Helper/loader/lib"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var Argfunc = []string{"__p___argv", "__p___argc", "GetCommandLineA", "GetCommandLineW", "__wgetmainargs", "__getmainargs"}

type ArgsFunc struct {
	Name    string
	Address uintptr
}

type ArgInjector func(addr uintptr, argv []string) error

var ArgInjectors = map[string]ArgInjector{
	"__p___argv": func(addr uintptr, argv []string) error {
		return InjectArgv(addr, argv)
	},
	"__p___argc": func(addr uintptr, argv []string) error {
		return InjectArgc(addr, argv)
	},
	"GetCommandLineA": func(addr uintptr, argv []string) error {
		return InjectCommandLineA(addr, argv)
	},
	"GetCommandLineW": func(addr uintptr, argv []string) error {
		return InjectCommandLineW(addr, argv)
	},
	"__wgetmainargs": func(addr uintptr, argv []string) error {
		return InjectCmdLn(addr, argv)
	},
	"__getmainargs": func(addr uintptr, argv []string) error {
		return InjectCmdLn(addr, argv)
	},
}

func UpdateExecMemory(funcAddr uintptr, sc []byte) (err error) {
	var empty uint32
	if err = windows.VirtualProtect(funcAddr, uintptr(len(sc)), 0x04, &empty); err != nil {
		return err
	}
	lib.Memcpy(uintptr(unsafe.Pointer(&sc[0])), funcAddr, uintptr(len(sc)))
	if err = windows.VirtualProtect(funcAddr, uintptr(len(sc)), 0x20, &empty); err != nil {
		return err
	}
	return nil
}

func InjectArgv(addr uintptr, argv []string) (err error) {

	var sc []byte
	ptrArgs := buildArgvPointers(argv)

	// movabs rax, entrypoint
	// ret
	opcode := fmt.Sprintf("48b8%xc3", formatPtr(ptrArgs))
	if sc, err = hex.DecodeString(opcode); err != nil {
		return err
	}

	return UpdateExecMemory(addr, sc)
}

func buildArgvPointers(argvs []string) unsafe.Pointer {
	ptrAddrAllArgs := make([]byte, 0)
	addrAllArgs := make([]byte, 0)

	for _, s := range argvs {
		strPtr := createStrPtr(s)
		addrAllArgs = append(addrAllArgs, formatPtr(strPtr)...)
	}
	addrAllArgs = append(addrAllArgs, formatAddr(0x0000000000000000)...)
	ptrAddrAllArgs = append(ptrAddrAllArgs, formatPtr(unsafe.Pointer(&addrAllArgs[0]))...)
	return unsafe.Pointer(&ptrAddrAllArgs[0])
}

func createStrPtr(str string) unsafe.Pointer {
	strBytes := make([]byte, 0)
	strBytes = append(strBytes, []byte(str)...)
	strBytes = append(strBytes, 0x00)
	return unsafe.Pointer(&strBytes[0])
}

func formatPtr(ptr unsafe.Pointer) []byte {
	return formatAddr(lib.PtrValue(ptr))
}

func formatAddr(addr uintptr) []byte {
	size := unsafe.Sizeof(uintptr(0))
	b := make([]byte, size)
	switch size {
	case 4:
		binary.LittleEndian.PutUint32(b, uint32(addr))
	default:
		binary.LittleEndian.PutUint64(b, uint64(addr))
	}
	return b
}

func InjectArgc(addr uintptr, argv []string) (err error) {
	var sc []byte
	argc := len(argv)
	argcBytes := formatAddr(uintptr(argc))
	addrArgc := unsafe.Pointer(&argcBytes[0])

	// movabs rax, entrypoint
	// ret
	opcode := fmt.Sprintf("48b8%xc3", formatPtr(addrArgc))
	if sc, err = hex.DecodeString(opcode); err != nil {
		return err
	}

	return UpdateExecMemory(addr, sc)
}

func InjectCommandLineA(addr uintptr, argv []string) (err error) {
	var sc []byte
	cmdLine := strings.Join(argv, " ")
	addrCmdLine := createStrPtr(cmdLine)

	// movabs rax, entrypoint
	// ret
	opcode := fmt.Sprintf("48b8%xc3", formatPtr(addrCmdLine))
	if sc, err = hex.DecodeString(opcode); err != nil {
		return err
	}

	return UpdateExecMemory(addr, sc)
}

func InjectCommandLineW(addr uintptr, argv []string) (err error) {
	var sc []byte
	cmdLine := strings.Join(argv, " ")
	//runes := utf16.Encode([]rune(cmdLine))
	//addrCmdLine := unsafe.Pointer(&runes[0])
	runes, _ := syscall.UTF16PtrFromString(cmdLine)
	addrCmdLine := unsafe.Pointer(runes)
	// movabs rax, entrypoint
	// ret
	opcode := fmt.Sprintf("48b8%xc3", formatPtr(addrCmdLine))
	if sc, err = hex.DecodeString(opcode); err != nil {
		return err
	}

	return UpdateExecMemory(addr, sc)
}

func InjectCmdLn(addr uintptr, argv []string) (err error) {
	var wCmdLine, aCmdLine uintptr

	msvcrtDLL, err := syscall.LoadLibrary("msvcrt.dll")
	if err != nil {
		return err
	}
	if wCmdLine, err = syscall.GetProcAddress(msvcrtDLL, "_wcmdln"); err != nil {
		return err
	}
	if aCmdLine, err = syscall.GetProcAddress(msvcrtDLL, "_acmdln"); err != nil {
		return err
	}

	cmdLine := strings.Join(argv, " ")
	addrCmdLine := createStrPtr(cmdLine)

	//runes := utf16.Encode([]rune(cmdLine))
	runes, _ := syscall.UTF16FromString(cmdLine)
	runes = append(runes, 0x00)
	addrCmdLineUnicode := unsafe.Pointer(&runes[0])

	var empty uint32

	if err = windows.VirtualProtect(wCmdLine, unsafe.Sizeof(uintptr(0)), 0x04, &empty); err != nil {
		return err
	}
	if err = windows.VirtualProtect(aCmdLine, unsafe.Sizeof(uintptr(0)), 0x04, &empty); err != nil {
		return err
	}

	lib.Memcpy(lib.PtrValue(unsafe.Pointer(&addrCmdLineUnicode)), wCmdLine, unsafe.Sizeof(uintptr(0)))
	lib.Memcpy(lib.PtrValue(unsafe.Pointer(&addrCmdLine)), aCmdLine, unsafe.Sizeof(uintptr(0)))

	return err
}
