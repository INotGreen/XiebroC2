package Proxy

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/Acebond/ReverseSocks5/statute"
)

// ServeConn is used to serve a single connection.
func ServeConn(conn net.Conn) error {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	// The client request detail
	request, err := ParseRequest(bufConn)
	if err != nil {
		if errors.Is(err, statute.ErrUnrecognizedAddrType) {
			if err := SendReply(conn, statute.RepAddrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("failed to send reply %w", err)
			}
		}
		return fmt.Errorf("failed to read destination address, %w", err)
	}

	if request.Request.Command != statute.CommandConnect &&
		request.Request.Command != statute.CommandBind &&
		request.Request.Command != statute.CommandAssociate {
		if err := SendReply(conn, statute.RepCommandNotSupported, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("unrecognized command[%d]", request.Request.Command)
	}

	//request.AuthContext = authContext
	request.LocalAddr = conn.LocalAddr()
	request.RemoteAddr = conn.RemoteAddr()
	// Process the client request
	return handleRequest(conn, request)
}
