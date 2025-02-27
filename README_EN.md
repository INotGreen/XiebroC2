<p align="center">
    <img src="https://badgen.net/github/stars/INotGreen/XiebroC2/?icon=github&color=black">
    <a href="https://github.com/INotGreen/XiebroC2/releases"><img src="https://img.shields.io/github/downloads/INotGreen/XiebroC2/total?color=blueviolet"></a>
    <img src="https://badgen.net/github/issues/INotGreen/XiebroC2">
    <a href="https://github.com/INotGreen/XiebroC2/wiki" style="text-decoration:none;">
     <img src="https://img.shields.io/badge/Documentation-wiki-yellow">
    </a>
</p>

## Main Features

- **Implant**: Written in Golang, compatible with Windows, Linux, and MacOS (support for mobile platforms under consideration for future updates).
- **Teamserver**: Built with .NET 6.0, does not require the .NET Core runtime environment.
- **Controller**: Supports reverse shell, file management, process management, network traffic monitoring, memory loading, reverse proxy (based on the [IOX](https://github.com/EddieIvan01/iox) model), screenshots, process injection and migration, AV/EDR detection, inline PowerShell commands.
- **Memory Operations**: Supports loading PE files into memory on Windows/Linux, process injection and migration, allowing file-free execution.
- **.NET Assemblies**: Execute .NET assemblies in memory (execute-assembly, inline-assembly).
- **Lua Scripting**: Extend command centers and menus through Lua scripts (similar to CNA scripts).
- **Custom RDI Shellcode**: (64-bit only, 32-bit requires manual client compilation) or generate shellcode using [donut](https://github.com/TheWover/donut) or [Godonut](https://github.com/Binject/go-donut).
- **Telegram Integration**: Set up Telegram notifications for host check-ins by modifying the `profile.json` parameters for Chat ID and API Token.

## Supported Platforms

**Implant (Session)**

- **Windows**: Windows 7–11, Windows Server 2008–2022
- **Linux**: Supports glibc 2.17+ (e.g., Ubuntu, Debian, CentOS)
- **MacOS**: macOS 10.15+

The project is compiled using Go 1.20 for compatibility. Note that Go 1.20+ does not support Windows 7, Windows Server 2008, and some older Linux systems. The payload in XiebroC2 only supports x64 architecture. For older systems, you must compile the source code with Go versions 1.19–1.16.

**Teamserver**

- **Windows**: Windows 8–11, Windows Server 2012–2022
- **Linux**: Supports glibc 2.17+ systems.

## Screenshots

Topology Structure

![image-20250114152703571](Image/image-20250114152703571.png)

Command List

![image-20250114162852363](Image/image-20250114162852363.png)

Memory Loading Mimikatz

![image-20250114162708390](Image/image-20250114162708390.png)

File Management

![image-20250114162940873](Image/image-20250114162940873.png)

Reverse Proxy

![image-20250114180254731](Image/image-20250114180254731.png)

## How to Use

- Download binaries directly from: [Release](https://github.com/INotGreen/XiebroC2/releases)
- Usage Guide: [XiebroC2 Wiki](https://github.com/INotGreen/XiebroC2/wiki)
- Extend penetration testing tools into Lua plugins: [Xiebro-Plugins](https://github.com/INotGreen/Xiebro-Plugins)

## Network Topology

View network traffic distribution with a visual topology diagram.

![Network Topology](Image/image-20240818150942815.png)

## Video Demo

[Demo](https://private-user-images.githubusercontent.com/89376703/305162512-771c2e88-afd8-493d-a575-7e10149837dd.mp4)

## TODO

- Develop payloads for PowerShell, VBScript, HTA, JScript, etc.
- Open more forms and API interfaces to facilitate Lua plugin development.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=INotGreen/XiebroC2&type=Date)](https://star-history.com/#INotGreen/XiebroC2&Date)

## Disclaimer

This project is intended solely for educational and research purposes in penetration testing practice. It is currently in a testing phase. It is strictly prohibited to use this tool for any illegal activities, including black market operations or unauthorized penetration attempts. The internet is not a lawless space! By using this tool, you agree to comply with these terms.

To prevent misuse by malicious actors, the most harmful features have been removed, leaving only basic functions for penetration testing demonstrations. The **Teamserver** and **Controller** components are not open-source.
