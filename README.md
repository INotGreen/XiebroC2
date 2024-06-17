

## Features

- The client side (Implant) is written in Golang, compatible with Windows, Linux, MacOS (mobile platforms are under consideration for future updates).
- The team server (Teamserver) is written in .net 8.0 and AOT compiled, featuring low memory usage without the need for any dependencies, nearly compatible with all platform systems.
- The controller supports reverse shell, file management, process management, network traffic monitoring, memory loading, reverse proxy (based on [IOX](https://github.com/EddieIvan01/iox) model), and screenshots.
- Supports in-memory loading of PE files on Windows/Linux, allowing the execution of trojans without dropping files to disk, and facilitating the use of third-party C2/RATs.
- Supports in-memory execution of .net assemblies (execute-assembly, inline-assembly).
- Supports extension of UI widgets, Session commands, and payload generation through lua (similar to CobaltStrike's cna scripts).
- Custom RDI shellcode support (64-bit only, 32-bit requires manual client compilation) or use [donut](https://github.com/TheWover/donut), [Godonut](https://github.com/Binject/go-donut) to generate your own shellcode.
- Teamserver supports hosting binary files, text, pictures (similar to SimpleHttpServer).
- Customizable team server configuration files, with custom Telegram chat ID/Token for notifications.
- The Controller UI is lightweight, with memory usage approximately 1/60th of CobaltStrike and 1/10th of Metasploit.
- Golang's compiler features have been blacklisted by some AV/EDR manufacturers, resulting in poor evasion capabilities.

## Supported Platforms

**Implant(Session)**

|    Windows (x86_x64)     | Linux (x86_x64) | MacOS |
| :----------------------: | :-------------: | :---: |
|        Windows 11        |     Ubuntu      | AMD64 |
|        Windows 10        |     Debian      | i386  |
|      Windows 8/8.1       |     CentOS      |  M1   |
|        Windows 7         |     ppc64le     |  M2   |
|        Windows XP        |      mips       |       |
| Windows Server 2000-2022 |      s390x      |       |

The payload in XiebroC2 currently only supports the x64-bit AMD architecture. If you have application scenarios in other environments, you need to compile the Go source code yourself.

## How to use

[xiebroC2 instruction manual](https://github.com/INotGreen/XiebroC2/wiki)

Write simply  plugins：[Xiebro-Plugins](https://github.com/INotGreen/Xiebro-Plugins)



## Topology

See network traffic distribution through a visual topology map

![image](https://github.com/INotGreen/XiebroC2/blob/main/Image/image-20240616214300666.png)



## TODO

- Currently, only the Session mode of the TCP/WebSocket protocol is supported. They are replacements for https. We will consider developing a reliable UDP protocol and support the Beacon mode in the future.
- Consider developing Powershell, VBscript, Hta, Jscript and other payloads.
- Open more forms and API interfaces to facilitate Lua extension plug-ins



## Disclaimer

This project is intended for educational and research purposes only in penetration testing exercises and is in beta. It is prohibited to use it for any illegal activities (including black market transactions, unauthorized penetration attacks, or financial exploitation)! The Internet is not a lawless space! If you choose to use this tool, you must comply with the above requirements.

In order to prevent the tool from being used by criminals, I have deleted the most harmful functions and only left some functions as penetration test drill demos. Teamserver and Controller are not open source.

