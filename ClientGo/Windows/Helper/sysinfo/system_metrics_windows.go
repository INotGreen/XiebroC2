package sysinfo

import (
	"syscall"
)

func GetSystemMetrics(index int) (int, error) {
	user32 := syscall.NewLazyDLL("User32.dll")
	procGetSystemMetrics := user32.NewProc("GetSystemMetrics")
	v, _, err := procGetSystemMetrics.Call(uintptr(index))
	// fmt.Printf("%#v \n", v)
	// fmt.Println(err)
	return int(v), err
}
