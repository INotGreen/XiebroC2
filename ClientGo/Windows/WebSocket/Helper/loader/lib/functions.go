package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func NtH(baseAddress uintptr) *IMAGE_NT_HEADERS {
	return (*IMAGE_NT_HEADERS)(unsafe.Pointer(baseAddress + uintptr((*IMAGE_DOS_HEADER)(unsafe.Pointer(baseAddress)).E_lfanew)))
}

func PtrOffset(ptr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(PtrValue(ptr) + offset)
}

func GetRelocTable(ntHeader *IMAGE_NT_HEADERS) *IMAGE_DATA_DIRECTORY {
	returnTable := &ntHeader.OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_BASERELOC]
	if returnTable.VirtualAddress == 0 {
		return nil
	} else {
		return returnTable
	}
}

func IsMSBSet(num uint) bool {
	uintSize := 32 << (^uint(0) >> 32 & 1)
	return num>>(uintSize-1) == 1
}

func ParseOrdinal(ordinal uint) (unsafe.Pointer, string) {
	funcOrdinal := uint16(ordinal)
	ptrName := unsafe.Pointer(uintptr(funcOrdinal))
	funcName := fmt.Sprintf("#%d", funcOrdinal)
	return ptrName, funcName
}

func CstrVal(ptr unsafe.Pointer) (out []byte) {
	var byteVal byte
	out = make([]byte, 0)
	for i := 0; ; i++ {
		byteVal = *(*byte)(unsafe.Pointer(ptr))
		if byteVal == 0x00 {
			break
		}
		out = append(out, byteVal)
		ptr = PtrOffset(ptr, 1)
	}
	return out
}

func PtrValue(ptr unsafe.Pointer) uintptr {
	return uintptr(unsafe.Pointer(ptr))
}

func ParseFuncAddress(base, offset uintptr) (unsafe.Pointer, string) {
	pImageImportByName := (*ImageImportByName)(unsafe.Pointer(base + offset))
	ptrName := unsafe.Pointer(&pImageImportByName.Name)
	funcName := string(CstrVal(ptrName))
	return ptrName, funcName
}

func GetProcAddress(libraryAddress, ptrName unsafe.Pointer) (uintptr, error) {
	getProcAddress := syscall.NewLazyDLL("kernel32").NewProc("GetProcAddress")
	ret, _, err := getProcAddress.Call(
		PtrValue(libraryAddress),
		PtrValue(ptrName))

	if err != syscall.Errno(0) {
		return 0, err
	}
	return ret, nil
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Memcpy(src, dst, size uintptr) {
	for i := uintptr(0); i < size; i++ {
		*(*byte)(unsafe.Pointer(dst + i)) = *(*byte)(unsafe.Pointer(src + i))
	}
}

func Memset(ptr uintptr, c byte, n uintptr) {
	var i uintptr
	for i = 0; i < n; i++ {
		pByte := (*byte)(unsafe.Pointer(ptr + i))
		*pByte = c
	}
}

func ObfuscateStrings(b []byte, blacklist []string) []byte {
	fmt.Printf("Replapcing %d keywords...\n", len(blacklist))
	for _, word := range blacklist {
		b = ReplaceWord(b, word)
	}
	return b
}

func ReplaceWord(b []byte, word string) []byte {
	newWord := shuffle(word)
	re := regexp.MustCompile("(?i)" + utf16LeStr(word))
	b = re.ReplaceAll(b, utf16Le(newWord))
	re2 := regexp.MustCompile("(?i)" + word)
	b = re2.ReplaceAll(b, []byte(newWord))

	return b
}

func utf16Le(s string) []byte {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	var buf bytes.Buffer
	t := transform.NewWriter(&buf, enc)
	t.Write([]byte(s))
	return buf.Bytes()
}

func utf16LeStr(s string) string {
	return string(utf16Le(s))
}

func shuffle(in string) string {
	rand.Seed(time.Now().Unix())
	inRune := []rune(in)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func DecodeB64(message string) string {
	decoded, _ := base64.StdEncoding.DecodeString(message)
	return string(decoded)
}
