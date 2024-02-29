package sysinfo

import (
	"syscall"
	"unsafe"
)

/*
Windows NT 4
Windows 95
Windows 98
Windows Me
Windows 2000
Windows XP
Windows XP 64
Windows Server 2003
Windows Server 2003 R2
Windwos Vista
Windows Server 2008
Windwos 7
Windows Server 2008 R2
Windwos 8
Windows Server 2012
Windows 8.1
Windows Server 2012 R2
Windows 10
Windows Server 2016
Windows Server 2019
Windows 11
Windows 11 +
*/
func WindosVersion() string {

	version, err := OSVersion()
	if err != nil {
		return ""
	}
	return version
}
func OSVersion() (string, error) {
	majorVersion, minorVersion, buildNumber := RtlGetNtVersionNumbers()
	// fmt.Printf("majorVersion:%d ,minorVersion:%d ,buildNumber:%d \n", majorVersion, minorVersion, buildNumber)
	o, err := GetVersionExW()
	if err != nil {
		// fmt.Println("GetVersionExW : ", err)
		return "", err
	}

	if majorVersion > 6 || (majorVersion == 6 && minorVersion >= 3) { // win8plus
		if majorVersion == 6 && minorVersion >= 3 {
			// Win8.1       : 6.3.9600
			// Windows Server 2012
			if o.ProductType == byte(VER_NT_WORKSTATION) {
				return "Windows 8.1", nil
			} else {
				// fmt.Println("o.ProductType :", o.ProductType)
				return "Windows Server 2012 R2", nil
			}
		} else if majorVersion == 10 && minorVersion == 0 {

			if o.ProductType == byte(VER_NT_WORKSTATION) {
				if buildNumber >= 22000 {
					return "Windows 11", nil
				} else { // if buildNumber >= 18363 {
					return "Windows 10", nil
				}
			} else {
				if buildNumber >= 17763 {
					return "Windows Server 2019", nil
				} else if buildNumber >= 14393 {
					return "Windows Server 2016", nil
				}
			}
		} else {
			return "Windows 11 +", nil
		}
	}

	s, err := GetNativeSystemInfo()
	if err != nil {
		return "", nil
	}
	u := s.GetDummyStructName()
	switch o.MajorVersion {
	case 4:
		{
			switch o.MinorVersion {
			case 0:
				{
					if int(o.PlatformId) == VER_PLATFORM_WIN32_NT {
						return "Windows NT 4", nil
					} else if int(o.PlatformId) == VER_PLATFORM_WIN32_WINDOWS {
						return "Windows 95", nil
					}
				}
			case 10:
				{
					return "Windows 98", nil
				}
			case 90:
				{
					return "Windows Me", nil
				}
			}
		}
	case 5:
		{
			switch o.MinorVersion {
			case 0:
				{
					return "Windows 2000", nil
				}
			case 1:
				{
					return "Windows XP", nil
				}
			case 2:
				{
					r2, err := GetSystemMetrics(SM_SERVERR2)
					if err != nil {
						return "", err
					}
					if o.ProductType == byte(VER_NT_WORKSTATION) && u.IsWin64() {
						return "Windows XP 64", nil
					} else if r2 == 0 {
						return "Windows Server 2003", nil
					} else if r2 != 0 {
						return "Windows Server 2003 R2", nil
					}
				}
			}
		}
	case 6:
		{
			switch o.MinorVersion {
			case 0:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windwos Vista", nil
					} else {
						return "Windows Server 2008", nil
					}
				}
			case 1:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windows 7", nil
					} else {
						return "Windows Server 2008 R2", nil
					}
				}
			case 2:
				{
					if o.ProductType == byte(VER_NT_WORKSTATION) {
						return "Windwos 8", nil
					} else {
						return "Windows Server 2012", nil
					}
				}
			}
		}
	}
	return "windows", nil
}

// Dll: ntdll.dll
// RtlGetNtVersionNumbers
// 获取系统的版本号
/*
   HINSTANCE hinst = LoadLibrary("ntdll.dll");
   DWORD dwMajor,dwMinor,dwBuildNumber;
   NTPROC proc = (NTPROC)GetProcAddress(hinst,"RtlGetNtVersionNumbers");
   proc(&dwMajor,&dwMinor,&dwBuildNumber);
   dwBuildNumber&=0xffff;
*/
func RtlGetNtVersionNumbers() (majorVersion, minorVersion, buildNumber uint32) {
	//var majorVersion, minorVersion, buildNumber uint32
	ntdll := syscall.NewLazyDLL("ntdll.dll")
	procRtlGetNtVersionNumbers := ntdll.NewProc("RtlGetNtVersionNumbers")
	//v, vv, err := procRtlGetNtVersionNumbers.Call(
	procRtlGetNtVersionNumbers.Call(
		uintptr(unsafe.Pointer(&majorVersion)),
		uintptr(unsafe.Pointer(&minorVersion)),
		uintptr(unsafe.Pointer(&buildNumber)),
	)
	// fmt.Printf("%#v \n", v)
	// fmt.Printf("%#v \n", vv)
	// fmt.Printf("%#v \n", err)
	// fmt.Println("开发版本:", buildNumber)
	buildNumber &= 0xffff
	// fmt.Println("-----------------RtlGetNtVersionNumbers-----------------")
	// fmt.Println("主版本号:", majorVersion)
	// fmt.Println("次版本号:", minorVersion)
	// fmt.Println("开发版本:", buildNumber)
	return
}
