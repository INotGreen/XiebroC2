package sysinfo

import (
	"syscall"
	"unsafe"
	// 这里我手动修改了下,请在您的项目中注意引用
)

// api dll
// Kernel32.dll
// Kernel32 := syscall.NewLazyDLL("kernel32.dll")
// 判断windwos操作系统是32位还是64位
// GetNativeSystemInfo function (sysinfoapi.h)
// docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-getnativesysteminfo
/*
void GetNativeSystemInfo(
  [out] LPSYSTEM_INFO lpSystemInfo
);
*/
func GetNativeSystemInfo() (*SystemInfo, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetNativeSystemInfo := kernel32.NewProc("GetNativeSystemInfo")
	systemInformation := &SystemInfo{}
	_, _, err := procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(systemInformation)))
	if err != nil && err.Error() == "The operation completed successfully." {
		err = nil
	}
	// fmt.Println(err)
	// fmt.Printf("%#v \n", systemInformation)
	// fmt.Printf("%#v \n", systemInformation.DummyUnionName)
	// dummyStructName := (*DummyStructName)(unsafe.Pointer(&systemInformation.DummyUnionName))
	// fmt.Printf("dummyStructName : %#v \n", dummyStructName)
	// fmt.Println("架构          : ", dummyStructName.Architecture())
	// fmt.Println("处理器数量    : ", systemInformation.NumberOfProcessors)
	// fmt.Println("处理器类型    : ", systemInformation.ProcessorType)
	return systemInformation, err
}

// 参数
// SYSTEM_INFO structure (sysinfoapi.h)
// docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
/*
// DUMMYUNIONNAME
/*
struct {
    WORD wProcessorArchitecture;
    WORD wReserved;
} DUMMYSTRUCTNAME;
} DUMMYUNIONNAME;
*/
// DummyStructName 虚拟结构名
type DummyStructName struct {
	/*
	   ProcessorArchitecture 值和说明
	   --------------------------------|--------------------------------|
	   Value                           |   Meaning
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_AMD64    |   x64 (AMD or Intel)
	   9                               |
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_ARM      |   ARM
	   5                               |
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_ARM64    |   ARM64
	   12                              |
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_IA64     |   Intel Itanium-based
	   6                               |
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_INTEL    |   x86
	   0                               |
	   --------------------------------|--------------------------------|
	   PROCESSOR_ARCHITECTURE_UNKNOWN  |   未知架构
	   0xffff                          |
	   --------------------------------|--------------------------------|
	*/
	ProcessorArchitecture WORD // 已安装操作系统的处理器体系结构。
	Reserved              WORD
}

// OemId
type OemId struct {
	OemId DWORD
}

/*
   When `go env -w GOARCH=amd64`: NumberOfProcessors:0x10000,ProcessorType:0x7e050006,(Error.不是太清楚这个什么原因):
   &version.SystemInfo{DummyUnionName:0x100000000009, PageSize:0x10000, MinimumApplicationAddress:0x7ffffffeffff, MaximumApplicationAddress:0xff, ActiveProcessorMask:0x21d800000008, NumberOfProcessors:0x10000, ProcessorType:0x7e050006, AllocationGranularity:0x0, ProcessorLevel:0x0, ProcessorRevision:0x0}
   When `go env -w GOARCH=386`  : NumberOfProcessors:0x8, ProcessorType:0x21d8 (correct 这个才是正确的):
   &version.SystemInfo{DummyUnionName:0x9, PageSize:0x1000, MinimumApplicationAddress:0x10000, MaximumApplicationAddress:0xfffeffff, ActiveProcessorMask:0xff, NumberOfProcessors:0x8, ProcessorType:0x21d8, AllocationGranularity:0x10000, ProcessorLevel:0x6, ProcessorRevision:0x7e05}
*/
/*
// 参数类型
typedef struct _SYSTEM_INFO {
    union {
      DWORD dwOemId;
      struct {
        WORD wProcessorArchitecture;
        WORD wReserved;
      } DUMMYSTRUCTNAME;
    } DUMMYUNIONNAME;
    DWORD     dwPageSize;
    LPVOID    lpMinimumApplicationAddress;
    LPVOID    lpMaximumApplicationAddress;
    DWORD_PTR dwActiveProcessorMask;
    DWORD     dwNumberOfProcessors;
    DWORD     dwProcessorType;
    DWORD     dwAllocationGranularity;
    WORD      wProcessorLevel;
    WORD      wProcessorRevision;
  } SYSTEM_INFO, *LPSYSTEM_INFO;
*/
type SystemInfo struct {
	DummyUnionName            uintptr   //DummyUnionName Or OemId ,one of this Use Same Memory.两个结构使用同一块内存
	PageSize                  DWORD     // 虚拟内存页的大小
	MinimumApplicationAddress LPVOID    // 应用程序和动态链接库（DLL）可访问的最低内存地址
	MaximumApplicationAddress LPVOID    // 应用程序和动态链接库（DLL）可访问的最高内存地址
	ActiveProcessorMask       DWORD_PTR // 表示配置到系统中的处理器集的掩码
	NumberOfProcessors        DWORD     // 处理器数量
	ProcessorType             DWORD     // 处理器类型
	AllocationGranularity     DWORD     // 虚拟内存的起始地址
	ProcessorLevel            WORD      // 依赖于体系结构的处理器级别。它只能用于显示目的。要确定处理器的功能集，请使用IsProcessorFeaturePresent函数。
	ProcessorRevision         WORD      // 依赖于体系结构的处理器版本
}

// 当前系统中的中央处理器的架构
func (d *DummyStructName) Architecture() string {
	switch d.ProcessorArchitecture {
	case 0:
		{
			return "x86" // 32位
		}
	case 5:
		{
			return "ARM" // 32位
		}
	case 6:
		{
			return "Itanium" //Intel  Itanium-based // Intel 奔腾架构  32位处理器
		}
	case 9:
		{
			return "x64" //  (AMD or Intel)  64位处理器
		}
	case 12:
		{
			return "ARM64" // 64位处理器
		}
	case 0xffff:
		{
			return "Unknow" // 未知
		}
	default:
		{
			return "Unknow" // 未知
		}
	}
}

// 是否为64位操作系统
func (d *DummyStructName) IsWin64() bool {
	return d.ProcessorArchitecture == WORD(PROCESSOR_ARCHITECTURE_AMD64) || d.ProcessorArchitecture == WORD(PROCESSOR_ARCHITECTURE_IA64)
}

// 处理器类型
// PROCESSOR_INTEL_386      (386)
// PROCESSOR_INTEL_486      (486)
// PROCESSOR_INTEL_PENTIUM  (586)
// PROCESSOR_INTEL_IA64     (2200)
// PROCESSOR_AMD_X8664      (8664)
// PROCESSOR_ARM            (Reserved)
func (s *SystemInfo) GetProcessorType() string {
	switch s.ProcessorType {
	case 386: // 0x0182
		{
			return "386"
		}
	case 486: // 0x01E6
		{
			return "486"
		}
	case 586: // 0x024A
		{
			return "Pentium" // "奔腾"
		}
	case 2200: // 0x0898
		{
			return "Itanium" //"安腾"
		}
	case 8664: // 0x21D8
		{
			return "X8664" //
		}
		// case xxx{ return "ARM" }
	default:
		{
			return "Unknow"
		}
	}
}
func (s *SystemInfo) GetDummyStructName() *DummyStructName {
	dummyStructName := (*DummyStructName)(unsafe.Pointer(&s.DummyUnionName))
	// fmt.Printf("dummyStructName : %#v \n", dummyStructName)
	// fmt.Println("架构          : ", dummyStructName.Architecture())
	// fmt.Println("处理器数量    : ", s.NumberOfProcessors)
	// fmt.Println("处理器类型    : ", s.ProcessorType)
	return dummyStructName
}
