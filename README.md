



## 使用说明（英文）

[XiebroC2-v3.1-README](https://github.com/INotGreen/XiebroC2/blob/main/README_EN.md)



## 特点/特征

- 被控端(Client)由Golang编写，兼容WIndows、Linux、MacOS上线（未来会考虑移动端上线）

- 团队服务器(Teamserver)由.net 8.0 编写、AOT编译，内存占用低，无需安装任何依赖，几乎可以兼容全平台系统

- 控制端(Controller)支持反弹shell，文件管理、进程管理、网络流量监控、内存加载、自定义UI背景色等功能

- 支持无文件落地，内存执行shellcode、.NET 程序集（execute-assembly、inline-assembly)、PE文件（如内存加载fscan等扫描器、POC/EXP)

- 支持websocket模式，以及cdn、域名上线PC

- 支持反向代理功能，类似于frps、ew、Stowaway、并且速度不逊色于它们

- 支持通过lua扩展UI控件、Session命令和载荷生成（类似于CobaltStrike的cna脚本）

- 支持Teamserver托管二进制文件、文本、图片(类似SimpleHttpServer)

- 支持Teamserver自定义配置文件：自定义内存加载方式（Fork&&run 或者Inline），自定义前置rdiShellcode64（仅64位，32位需要手动编译client）、Telegram的chat ID/Token上线通知、Websocket路由特征。

- 控制端(Controller)UI轻量级交互界面，内存占用大约是CobaltStrike的60分之一，是Metasploit的10分之一

- 与Beacon模式不同的是，被控端是Session模式，可以用netstat查看实时连接端口，并且流量通信也是实时性的

- 由于Golang编译器的代码结构比较复杂，杀毒软件很难对Go的二进制文件进行准确的静态分析，随着时间的推移，Golang被越来越多的AV/EDR厂商标记为恶意软件其中包括（360、微软、Google、Elastic、Ikarus）

  

## 支持平台

**Client(Session)**

|    Windows（x86_x64）    | Linux(x86_x64) | MacOS |
| :----------------------: | :------------: | :---: |
|        Windows11         |     ubuntu     | AMD64 |
|        Windows10         |     Debian     | i386  |
|       Windows8/8.1       |     CentOS     |  M1   |
|         Windows7         |    ppc64le     |  M2   |
|        Windows-XP        |      mips      |       |
| Windows Server 2000-2022 |     s390x      |       |



## 快速使用

[快速使用](https://github.com/INotGreen/XiebroC2/wiki)

编写简单的插件：[插件编写](https://github.com/INotGreen/Xiebro-Plugins)



## 免责声明

本项目仅用于渗透测试演练的学习交流和研究，强烈不建议您用于任何的实际途径（包括黑灰产交易、非法渗透攻击、割韭菜），网络不是法外之地！如果您使用该工具则应该自觉遵守以上要求。

为了避免该工具被非法分子利用，所以本人已经将危害较大的功能删除，只留下部分功能作为渗透测试演练demo，teamserver和Controller不进行开源




