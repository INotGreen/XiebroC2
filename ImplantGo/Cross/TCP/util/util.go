package util

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func GetExecPath() (string, error) {
	pid := int32(os.Getpid())
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}
	for _, p := range processes {
		if p.Pid == pid {
			return p.Cmdline()
		}
	}
	return "", err

}

func GetExecPathEx() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func ReadMySelf() []byte {
	path, _ := GetExecPathEx()
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return f
}
