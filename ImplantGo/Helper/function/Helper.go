//go:build windows
// +build windows

package Function

import (
	"golang.org/x/sys/windows/registry"
)

func GetClrVersion() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full`, registry.QUERY_VALUE)
	if err != nil {
		return "v2" // If the registry cannot be accessed, assume CLR 2.0 is returned
	}
	defer key.Close()

	// If the Release key is present, CLR 4.0 or higher is installed
	if _, _, err := key.GetIntegerValue("Release"); err == nil {
		return "v4"
	}

	return "v2"
}
