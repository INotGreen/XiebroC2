//go:build windows
// +build windows

package syscalls

import (
	"main/Helper/sysinfo"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
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

func GetWindowsVersion() string {
	return sysinfo.WindosVersion()
}

var (
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
