



## 使用说明（英文）

[XiebroC2-v3.1-README](https://github.com/INotGreen/XiebroC2/blob/main/README_EN.md)

## 更新

[+] 2024.2.18 XiebroC2-3.1发布

[+] 2024.2.20 XiebroC2-3.1 修复bug： [xiebroc2](https://github.com/INotGreen/XiebroC2/releases/download/XieBroC2-v3.1/XiebroC2-v3.1.7z)

随缘更新中。。。

如果您喜欢该项目的话，可以左上角star + fork + follow，非常感谢！

## 特点

- 被控端(Client)由Golang编写，兼容WIndows、Linux、MacOS上线（未来会考虑移动端上线）

- 团队服务器(Teamserver)由.net 8.0 编写、AOT编译，内存占用低，无需安装任何依赖，几乎可以兼容全平台系统

- 控制端(Controller)支持反弹shell，文件管理、进程管理、网络流量监控、内存加载、自定义UI背景色等功能

- 支持Windows/Linux内存加载PE文件，即文件不落地执行木马，中转第三方C2/RAT

- 支持内存执行.net 程序集（execute-assembly、inline-assembly)

- 支持通过lua扩展UI控件、Session命令，载荷生成（类似于CobaltStrike的cna脚本）

- 支持自定义RDIshellcode（仅64位，32位需要手动编译client）或者用[donut](https://github.com/TheWover/donut)、[Godonut](https://github.com/Binject/go-donut)生成属于自己的shellcode

- 支持Teamserver托管二进制文件、文本、图片(类似SimpleHttpServer)

- 支持团队服务器自定义配置文件,自定义Telegram的chat ID/Token上线通知

- 控制端(Controller)UI轻量级交互界面，内存占用大约是CobaltStrike的60分之一，是Metasploit的10分之一

- Golang的编译器特征已经被部分AV/EDR厂商标黑了,因此免杀效果较差

  

  

## 支持平台

**Client(Session)**

|    Windows（x86_x64）    | Linux(x86_x64) | MacOS |
| :----------------------: | :------------: | :---: |
|        Windows11         |     ubuntu     | AMD64 |
|        Windows10         |     Debian     | i386  |
|       Windows8/8.1       |     CentOS     |  M1   |
|         Windows7         |    ppc64le     |  M2   |
|        Windwos-XP        |      mips      |       |
| Windows Server 2000-2022 |     s390x      |       |



## 快速使用

- 通过curl下载，密码：123456

```bash
curl -o XiebroC2-v3.1.7z https://github.com/INotGreen/XiebroC2/releases/download/XieBroC2-v3.1/Xiebro-v3.1.7z
```

- 控制端需要运行在.Net Framework4.8以上（Win10/11无需安装，win7需要下载: [.net framworkd4.8下载](https://dotnet.microsoft.com/zh-cn/download/dotnet-framework/thank-you/net48-offline-installer)）

- 修改TeamServerIP和TeamServerPort为VPS的IP和端口，然后保存为profile.json

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
    "rdiShellcode64": "",
    "rdiShellcode32": ""
}
```

服务端：

```bash
Teamserver.exe -c profile.json
```



## 上线演示

<video src="https://private-user-images.githubusercontent.com/89376703/305162512-771c2e88-afd8-493d-a575-7e10149837dd.mp4" width="640" height="480" controls></video>





## 命令列表




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





## 添加插件

<video src="https://private-user-images.githubusercontent.com/89376703/305687743-fb39df88-0f29-4359-9cd4-fc4bfa698270.mp4" width="640" height="480" controls></video>

## 拖拽式批量上传文件

<video src="https://private-user-images.githubusercontent.com/89376703/306153487-551e96db-9253-4a9f-8c2d-5c99c0280c8a.mp4" width="640" height="480" controls></video>



## 插件编写

- 学习编写lua插件:[Xiebro-Plugins](https://github.com/INotGreen/Xiebro-Plugins?tab=readme-ov-file#executeassembly)



## 计划进行

- 用C/C++/C#/Rust编写多阶段加载器（Multi-stage loader），体积尽量控制在150kb以内。

- 目前正反向代理和端口转发未开放，未来考虑完善和开发这个功能。

- 正在开发WebSocket/RUDP/SMB协议的Session模式，Beacon模式仅考虑开发HTTP/HTTPS/DNS。
- 考虑开发Powershell、VBscript、Hta、Jscript等载荷。

- 开放更多窗体和API接口，以便lua扩展插件

- 目前Golang的控制台隐藏问题还无法得到很好的方案，如果您知道如何解决这个问题请联系我。





## 免责声明

本项目仅用于渗透测试演练的学习交流和研究，强烈不建议您用于任何的实际途径（包括黑灰产交易、非法渗透攻击、割韭菜），网络不是法外之地！如果您使用该工具则应该自觉遵守以上要求



## 提交Bug和建议

<img src="Image\Image.jpg"  />



