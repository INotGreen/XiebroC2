package sysinfo

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type HANDLE uintptr
type HMODULE uintptr
type LPVOID uintptr
type WORD uint16
type DWORD uint32
type DWORD_PTR uintptr // 不确定
type ULong uint32
type HRESULT uint32 //
type LPBYTE *byte
type LPDWORD *uint32
type LPWSTR *uint16
type LMSTR *uint16

// type WCHAR = wchar_t
// type wchar_t = uint16
// UTF16toString converts a pointer to a UTF16 string into a Go string.
func UTF16toString(p *uint16) string {
	return syscall.UTF16ToString((*[4096]uint16)(unsafe.Pointer(p))[:])
}
func StrPtr(s string) uintptr {
	systemNameUTF16String, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		// 这里转换的时候出错,则不继续执行 OR 赋值用本地的
		tmp := utf16.Encode([]rune("\x00"))
		systemNameUTF16String = &tmp[0]
	}
	return uintptr(unsafe.Pointer(systemNameUTF16String))
}
func CharsPtr(s string) uintptr {
	bPtr, err := syscall.BytePtrFromString(s)
	if err != nil {
		return uintptr(0) // 这么写肯定不太对 @TODO
	}
	return uintptr(unsafe.Pointer(bPtr))
}
func IntPtr(n int) uintptr {
	return uintptr(n)
}

// winnt.go

// 判断操作系统的量
var (
	PROCESSOR_ARCHITECTURE_INTEL          = 0
	PROCESSOR_ARCHITECTURE_MIPS           = 1
	PROCESSOR_ARCHITECTURE_ALPHA          = 2
	PROCESSOR_ARCHITECTURE_PPC            = 3
	PROCESSOR_ARCHITECTURE_SHX            = 4
	PROCESSOR_ARCHITECTURE_ARM            = 5
	PROCESSOR_ARCHITECTURE_IA64           = 6
	PROCESSOR_ARCHITECTURE_ALPHA64        = 7
	PROCESSOR_ARCHITECTURE_MSIL           = 8
	PROCESSOR_ARCHITECTURE_AMD64          = 9
	PROCESSOR_ARCHITECTURE_IA32_ON_WIN64  = 10
	PROCESSOR_ARCHITECTURE_NEUTRAL        = 11
	PROCESSOR_ARCHITECTURE_ARM64          = 12
	PROCESSOR_ARCHITECTURE_ARM32_ON_WIN64 = 13
	PROCESSOR_ARCHITECTURE_IA32_ON_ARM64  = 14
	PROCESSOR_ARCHITECTURE_UNKNOWN        = 0xFFFF
	PRODUCT_UNDEFINED                     = 0x00000000
	PRODUCT_ULTIMATE                      = 0x00000001
	PRODUCT_HOME_BASIC                    = 0x00000002
	PRODUCT_HOME_PREMIUM                  = 0x00000003
	PRODUCT_ENTERPRISE                    = 0x00000004
	PRODUCT_HOME_BASIC_N                  = 0x00000005
	PRODUCT_BUSINESS                      = 0x00000006
	PRODUCT_STANDARD_SERVER               = 0x00000007
	PRODUCT_DATACENTER_SERVER             = 0x00000008
	PRODUCT_SMALLBUSINESS_SERVER          = 0x00000009
	PRODUCT_ENTERPRISE_SERVER             = 0x0000000A
	PRODUCT_STARTER                       = 0x0000000B
	PRODUCT_DATACENTER_SERVER_CORE        = 0x0000000C
	PRODUCT_STANDARD_SERVER_CORE          = 0x0000000D
	PRODUCT_ENTERPRISE_SERVER_CORE        = 0x0000000E
	PRODUCT_ENTERPRISE_SERVER_IA64        = 0x0000000F
	PRODUCT_BUSINESS_N                    = 0x00000010
	PRODUCT_WEB_SERVER                    = 0x00000011
	PRODUCT_CLUSTER_SERVER                = 0x00000012
	PRODUCT_HOME_SERVER                   = 0x00000013
	PRODUCT_STORAGE_EXPRESS_SERVER        = 0x00000014
	PRODUCT_STORAGE_STANDARD_SERVER       = 0x00000015
	PRODUCT_STORAGE_WORKGROUP_SERVER      = 0x00000016
	PRODUCT_STORAGE_ENTERPRISE_SERVER     = 0x00000017
	PRODUCT_SERVER_FOR_SMALLBUSINESS      = 0x00000018
	PRODUCT_SMALLBUSINESS_SERVER_PREMIUM  = 0x00000019
	PRODUCT_UNLICENSED                    = uint64(0xABCDABCD)
	VER_PLATFORM_WIN32s                   = 0
	VER_PLATFORM_WIN32_WINDOWS            = 1
	VER_PLATFORM_WIN32_NT                 = 2
	VER_NT_WORKSTATION                    = 0x0000001
	VER_NT_DOMAIN_CONTROLLER              = 0x0000002
	VER_NT_SERVER                         = 0x0000003
	//#if(_WIN32_WINNT >= 0x0501)
	SM_TABLETPC    = 86
	SM_MEDIACENTER = 87
	SM_STARTER     = 88
	SM_SERVERR2    = 89
	//#endif /* _WIN32_WINNT >= 0x0501 */
)
