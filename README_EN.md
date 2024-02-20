## Updates

[+] Released XiebroC2-3.1 on 2024.2.18

[+] Fixed bug in XiebroC2-3.1 on 2024.2.20: [xiebroc2](https://github.com/INotGreen/XiebroC2/releases/download/XieBroC2-v3.1/XiebroC2-v3.1.7z)

Updates will be released as they come.

If you like this project, feel free to star + fork + follow on the top left corner. Thank you very much!

## Features

- The client side (Client) is written in Golang, compatible with Windows, Linux, MacOS (mobile platforms are under consideration for future updates).
- The team server (Teamserver) is written in .net 8.0 and AOT compiled, featuring low memory usage without the need for any dependencies, nearly compatible with all platform systems.
- The controller (Controller) supports reverse shell, file management, process management, network traffic monitoring, memory loading, custom UI background colors, and more.
- Supports in-memory loading of PE files on Windows/Linux, allowing the execution of trojans without dropping files to disk, and facilitating the use of third-party C2/RATs.
- Supports in-memory execution of .net assemblies (execute-assembly, inline-assembly).
- Supports extension of UI widgets, Session commands, and payload generation through lua (similar to CobaltStrike's cna scripts).
- Custom RDI shellcode support (64-bit only, 32-bit requires manual client compilation) or use [donut](https://github.com/TheWover/donut), [Godonut](https://github.com/Binject/go-donut) to generate your own shellcode.
- Teamserver supports hosting binary files, text, pictures (similar to SimpleHttpServer).
- Customizable team server configuration files, with custom Telegram chat ID/Token for notifications.
- The Controller UI is lightweight, with memory usage approximately 1/60th of CobaltStrike and 1/10th of Metasploit.
- Golang's compiler features have been blacklisted by some AV/EDR manufacturers, resulting in poor evasion capabilities.

## Supported Platforms

**Client (Session)**

|    Windows (x86_x64)     | Linux (x86_x64) | MacOS |
| :----------------------: | :-------------: | :---: |
|        Windows 11        |     Ubuntu      | AMD64 |
|        Windows 10        |     Debian      | i386  |
|      Windows 8/8.1       |     CentOS      |  M1   |
|        Windows 7         |     ppc64le     |  M2   |
|        Windows XP        |      mips       |       |
| Windows Server 2000-2022 |      s390x      |       |

## Quick Start

- Download with curl, password: 123456

```
bashCopy code
curl -o XiebroC2-v3.1.7z https://github.com/INotGreen/XiebroC2/releases/download/XieBroC2-v3.1/Xiebro-v3.1.7z
```

- The controller requires .Net Framework 4.8 or higher (no installation required for Win10/11, Win7 needs to download: [.net framework 4.8 download](https://dotnet.microsoft.com/zh-cn/download/dotnet-framework/thank-you/net48-offline-installer))
- Modify TeamServerIP and TeamServerPort to your VPS's IP and port, then save as profile.json

```json
{
    "TeamServerIP": "192.168.31.81",
    "TeamServerPort": "8880",
    "Password": "123456",
    "StagerPort": "4050",
    "Telegram_Token": "",
    "Telegram_chat_ID": "",
    "Fork": true,
    "Process64": "C:\\windows\\system32\\notepad.exe",
    "Process86": "C:\\Windows\\SysWOW64\\notepad.exe",
    "WebServers": [],
    "listeners": [],
    "s_Reflection_dll_x64": ""
}
```

Server side:

```bash
Teamserver.exe -c profile.json
```

## Demonstration

Online demo

<video src="https://private-user-images.githubusercontent.com/89376703/305162512-771c2e88-afd8-493d-a575-7e10149837dd.mp4" width="640" height="480" controls></video>



## Command list


| Commands         |               Usage                |                   Description                    |
| :--------------- | :--------------------------------: | :----------------------------------------------: |
| nps              |     nps  “powershell command”      |        Unmanaged run powershell in memory        |
| Inline-assembly  | inline-assembly  “FilePath” “args” |           Inline execute .net assembly           |
| execute-assembly | execute-assembly “FilePath” ”args” | Fork child process execute loader .net  assembly |
| runpe            |      runpe  “FilePath” “args”      |    loader C/C++ PE in the memory for Windwos     |
| shell            |        shell “cmd command”         |                 Execute  command                 |
| powershell       |  powershell “powershell command”   |            Execute powershell command            |
| checkAV          |              checkAV               |             Detect AV/EDR processes              |
| upload           |   upload “RemotePath” “FilePath”   |            Upload File to the target             |
| memfd            |      memfd “FilePath” “args”       |        PE loader in the memory for Linux         |
| help             |                help                |                View command list                 |
| cls              |                cls                 |                   Clear screen                   |



## Add plugins

<video src="https://private-user-images.githubusercontent.com/89376703/305687743-fb39df88-0f29-4359-9cd4-fc4bfa698270.mp4" width="640" height="480" controls></video>

## Drag File

<video src="https://private-user-images.githubusercontent.com/89376703/306153487-551e96db-9253-4a9f-8c2d-5c99c0280c8a.mp4" width="640" height="480" controls></video>



## Plugin writing

- Learn to write lua plugins: [Xiebro-Plugins](https://github.com/INotGreen/Xiebro-Plugins?tab=readme-ov-file#executeassembly)

## Ongoing Plans

- Development of a multi-stage loader in C/C++/C#/Rust with the aim to keep the size under 150kb.
- Currently, reverse proxy and port forwarding features are not available but are considered for future development.
- Development of Session modes for WebSocket/RUDP/SMB protocols is underway, with Beacon mode limited to HTTP/HTTPS/DNS protocols.
- Consideration for developing payloads for  Powershell, VBscript, Hta, Jscript, etc.
- Opening more forms and API interfaces for lua extension plugins.
- The issue of console hiding in Golang is yet to be resolved satisfactorily. If you have solutions, please contact us.

## Disclaimer

This project is intended solely for educational and research purposes within penetration testing exercises. It is strongly advised against using it for any illegal activities (including black-market transactions, unauthorized penetration attacks, or financial exploitation). The internet is not a lawless space! If you choose to use this tool, you must comply with the above requirements.

