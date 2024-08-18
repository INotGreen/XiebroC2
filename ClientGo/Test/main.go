package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func getClrVersion() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full`, registry.QUERY_VALUE)
	if err != nil {
		return "v2.0" // 如果无法访问注册表，假定返回 CLR 2.0
	}
	defer key.Close()

	release, _, err := key.GetIntegerValue("Release")
	if err != nil {
		return "v2.0" // 如果读取 Release 键失败，假定返回 CLR 2.0
	}

	// 根据 Release 的值判断具体版本
	switch {
	case release >= 528040:
		return "v4.8"
	case release >= 461808:
		return "v4.7.2"
	case release >= 461308:
		return "v4.7.1"
	case release >= 460798:
		return "v4.7"
	case release >= 394802:
		return "v4.6.2"
	case release >= 394254:
		return "v4.6.1"
	case release >= 393295:
		return "v4.6"
	case release >= 379893:
		return "v4.5.2"
	case release >= 378675:
		return "v4.5.1"
	case release >= 378389:
		return "v4.5"
	default:
		return "v4.0"
	}
}

func main() {
	clrVersion := getClrVersion()
	fmt.Println("CLR Version:", clrVersion)
}
