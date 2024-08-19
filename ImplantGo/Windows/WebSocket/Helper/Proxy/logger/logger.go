package logger

import (
	"fmt"
	"main/Helper/Proxy/option"
	"os"
)

const (
	pWARN    = "[!] "
	pINFO    = "[+] "
	pSUCCESS = "[*] "
)

func Info(format string, args ...interface{}) {
	if option.VERBOSE {
		fmt.Fprintf(os.Stdout, pINFO+format+"\n", args...)
	}
}

func Warn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, pWARN+format+"\n", args...)
}

func Success(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, pSUCCESS+format+"\n", args...)
}
