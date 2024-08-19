package lib

const (
	IMAGE_NUMBEROF_DIRECTORY_ENTRIES = 16
	IMAGE_DIRECTORY_ENTRY_IMPORT     = 1 // Import Directory
	IMAGE_DIRECTORY_ENTRY_BASERELOC  = 5 // Base Relocation Table
)

type (
	DWORD     uint32
	ULONGLONG uint64
	WORD      uint16
	BYTE      uint8
	LONG      uint32
)

type _IMAGE_DATA_DIRECTORY struct {
	VirtualAddress DWORD
	Size           DWORD
}
type IMAGE_DATA_DIRECTORY _IMAGE_DATA_DIRECTORY

type IMAGE_OPTIONAL_HEADER64 _IMAGE_OPTIONAL_HEADER64
type IMAGE_OPTIONAL_HEADER IMAGE_OPTIONAL_HEADER64

type _IMAGE_OPTIONAL_HEADER64 struct {
	Magic                       WORD
	MajorLinkerVersion          BYTE
	MinorLinkerVersion          BYTE
	SizeOfCode                  DWORD
	SizeOfInitializedData       DWORD
	SizeOfUninitializedData     DWORD
	AddressOfEntryPoint         DWORD
	BaseOfCode                  DWORD
	ImageBase                   ULONGLONG
	SectionAlignment            DWORD
	FileAlignment               DWORD
	MajorOperatingSystemVersion WORD
	MinorOperatingSystemVersion WORD
	MajorImageVersion           WORD
	MinorImageVersion           WORD
	MajorSubsystemVersion       WORD
	MinorSubsystemVersion       WORD
	Win32VersionValue           DWORD
	SizeOfImage                 DWORD
	SizeOfHeaders               DWORD
	CheckSum                    DWORD
	Subsystem                   WORD
	DllCharacteristics          WORD
	SizeOfStackReserve          ULONGLONG
	SizeOfStackCommit           ULONGLONG
	SizeOfHeapReserve           ULONGLONG
	SizeOfHeapCommit            ULONGLONG
	LoaderFlags                 DWORD
	NumberOfRvaAndSizes         DWORD
	DataDirectory               [IMAGE_NUMBEROF_DIRECTORY_ENTRIES]IMAGE_DATA_DIRECTORY
}

type _IMAGE_FILE_HEADER struct {
	Machine              WORD
	NumberOfSections     WORD
	TimeDateStamp        DWORD
	PointerToSymbolTable DWORD
	NumberOfSymbols      DWORD
	SizeOfOptionalHeader WORD
	Characteristics      WORD
}

type IMAGE_FILE_HEADER _IMAGE_FILE_HEADER

type _IMAGE_DOS_HEADER struct { // DOS .EXE header
	E_magic    WORD     // Magic number
	E_cblp     WORD     // Bytes on last page of file
	E_cp       WORD     // Pages in file
	E_crlc     WORD     // Relocations
	E_cparhdr  WORD     // Size of header in paragraphs
	E_minalloc WORD     // Minimum extra paragraphs needed
	E_maxalloc WORD     // Maximum extra paragraphs needed
	E_ss       WORD     // Initial (relative) SS value
	E_sp       WORD     // Initial SP value
	E_csum     WORD     // Checksum
	E_ip       WORD     // Initial IP value
	E_cs       WORD     // Initial (relative) CS value
	E_lfarlc   WORD     // File address of relocation table
	E_ovno     WORD     // Overlay number
	E_res      [4]WORD  // Reserved words
	E_oemid    WORD     // OEM identifier (for E_oeminfo)
	E_oeminfo  WORD     // OEM information; E_oemid specific
	E_res2     [10]WORD // Reserved words
	E_lfanew   LONG     // File address of new exe header
}

type IMAGE_DOS_HEADER _IMAGE_DOS_HEADER

type IMAGE_IMPORT_DESCRIPTOR struct {
	OriginalFirstThunk DWORD
	TimeDateStamp      DWORD
	ForwarderChain     DWORD
	Name               DWORD
	FirstThunk         DWORD
}

type ImageThunkData struct {
	AddressOfData uintptr
}

type OriginalImageThunkData struct {
	Ordinal uint
}

type ImageImportByName struct {
	Hint uint16
	Name byte
}

type IMAGE_BASE_RELOCATION struct {
	VirtualAddress DWORD
	SizeOfBlock    DWORD
	//  WORD    TypeOffset[1];
}

type ImageReloc struct {
	OffsetType uint16
}

func (c *ImageReloc) GetOffset() uint16 {
	return c.OffsetType & 0x0fff
}

func (c *ImageReloc) GetType() uint16 {
	return (c.OffsetType & 0xf000) >> 12
}

type _IMAGE_NT_HEADERS64 struct {
	Signature      DWORD
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER
}
type IMAGE_NT_HEADERS64 _IMAGE_NT_HEADERS64
type IMAGE_NT_HEADERS IMAGE_NT_HEADERS64
