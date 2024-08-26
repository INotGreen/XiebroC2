package HandlePacket

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func executeCommandAndHandleCD(cmdString string) {
	if strings.HasPrefix(cmdString, "cd ") {
		arg := strings.TrimSpace(cmdString[3:])

		if arg == ".." {

			ProcessPath = filepath.Dir(ProcessPath)
		} else {

			if dirExists, err := directoryExists(arg); !dirExists {
				if err != nil {
					log.Printf("检查目录存在时发生错误: %v\n", err)
				} else {
					log.Printf("目录不存在: %s\n", arg)
				}
			} else {
				log.Printf("目标目录: %s\n", arg)
			}
			ProcessPath = arg

		}
	}
}

// directoryExists 检查指定的路径是否存在且为目录
func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil //
		}
		return false, err
	}
	return info.IsDir(), nil
}