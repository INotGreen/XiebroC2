// code from https://github.com/ring04h/s5.go
package socks5

import (
	"errors"
	"io"
	"main/Helper/Proxy/logger"
	"main/Helper/Proxy/netio"
	"main/Helper/Proxy/option"
	"net"
	"strconv"
	"time"
)

var (
	Commands = []string{"CONNECT", "BIND", "UDP ASSOCIATE"}
	AddrType = []string{"", "IPv4", "", "Domain", "IPv6"}
	Verbose  = false

	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support noauth method")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks only support connect command")
)

const (
	socksVer5       = 0x05
	socksCmdConnect = 0x01
)

func readAtLeast(r netio.Ctx, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}

	for n < min && err == nil {
		var nn int
		nn, err = r.DecryptRead(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func handShake(conn netio.Ctx) (err error) {
	const (
		idVer     = 0
		idNmethod = 1
	)

	buf := make([]byte, 258)

	var n int

	// make sure we get the nmethod field
	if n, err = readAtLeast(conn, buf, idNmethod+1); err != nil {
		return
	}

	if buf[idVer] != socksVer5 {
		return errVer
	}

	nmethod := int(buf[idNmethod]) //  client support auth mode
	msgLen := nmethod + 2          //  auth msg length
	if n == msgLen {               // handshake done, common case
		// do nothing, jump directly to send confirmation
	} else if n < msgLen { // has more methods to read, rare case
		if _, err = readAtLeast(conn, buf[n:msgLen], len(buf[n:msgLen])); err != nil {
			return
		}
	} else { // error, should not get extra data
		return errAuthExtraData
	}
	/*
	   X'00' NO AUTHENTICATION REQUIRED
	   X'01' GSSAPI
	   X'02' USERNAME/PASSWORD
	   X'03' to X'7F' IANA ASSIGNED
	   X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	   X'FF' NO ACCEPTABLE METHODS
	*/
	// send confirmation: version 5, no authentication required
	_, err = conn.EncryptWrite([]byte{socksVer5, 0})

	return
}

func parseTarget(conn netio.Ctx) (host string, err error) {
	const (
		idVer   = 0
		idCmd   = 1
		idType  = 3 // address type index
		idIP0   = 4 // ip addres start index
		idDmLen = 4 // domain address length index
		idDm0   = 5 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4   = 3 + 1 + net.IPv4len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv4 + 2port
		lenIPv6   = 3 + 1 + net.IPv6len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv6 + 2port
		lenDmBase = 3 + 1 + 1 + 2           // 3 + 1addrType + 1addrLen + 2port, plus addrLen
	)
	// refer to getRequest in server.go for why set buffer size to 263
	buf := make([]byte, 263)
	var n int

	// read till we get possible domain length field
	if n, err = readAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}

	// check version and cmd
	if buf[idVer] != socksVer5 {
		err = errVer
		return
	}

	/*
	   CONNECT X'01'
	   BIND X'02'
	   UDP ASSOCIATE X'03'
	*/

	if buf[idCmd] > 0x03 || buf[idCmd] == 0x00 {
		logger.Info("Unknown Command: %d", buf[idCmd])
	}

	if buf[idCmd] != socksCmdConnect { //  only support CONNECT mode
		err = errCmd
		return
	}

	// read target address
	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm: // domain name
		reqLen = int(buf[idDmLen]) + lenDmBase
	default:
		err = errAddrType
		return
	}

	if n == reqLen {
		// common case, do nothing
	} else if n < reqLen { // rare case
		if _, err = readAtLeast(conn, buf[n:reqLen], len(buf[n:reqLen])); err != nil {
			return
		}
	} else {
		err = errReqExtraData
		return
	}

	switch buf[idType] {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
	}
	port := bigEndianUint16(buf[reqLen-2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))

	return
}

func bigEndianUint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func pipeWhenClose(conn netio.Ctx, target string) {
	remoteConn, err := net.DialTimeout(
		"tcp", target,
		time.Millisecond*time.Duration(option.TIMEOUT),
	)
	if err != nil {
		logger.Info("Connect remote :" + err.Error())
		return
	}
	defer remoteConn.Close()

	tcpAddr := remoteConn.LocalAddr().(*net.TCPAddr)
	if tcpAddr.Zone == "" {
		if tcpAddr.IP.Equal(tcpAddr.IP.To4()) {
			tcpAddr.Zone = "ip4"
		} else {
			tcpAddr.Zone = "ip6"
		}
	}

	rep := make([]byte, 256)
	rep[0] = 0x05
	rep[1] = 0x00 // success
	rep[2] = 0x00 //RSV

	//IP
	if tcpAddr.Zone == "ip6" {
		rep[3] = 0x04 //IPv6
	} else {
		rep[3] = 0x01 //IPv4
	}

	var ip net.IP
	if "ip6" == tcpAddr.Zone {
		ip = tcpAddr.IP.To16()
	} else {
		ip = tcpAddr.IP.To4()
	}
	pindex := 4
	for _, b := range ip {
		rep[pindex] = b
		pindex += 1
	}
	rep[pindex] = byte((tcpAddr.Port >> 8) & 0xff)
	rep[pindex+1] = byte(tcpAddr.Port & 0xff)

	conn.EncryptWrite(rep[0 : pindex+2])
	// Transfer data

	remoteConnCtx, err := netio.NewTCPCtx(remoteConn, false)
	if err != nil {
		logger.Info("Socks5 remote connect error: %s", err.Error())
		return
	}

	netio.PipeForward(conn, remoteConnCtx)
}

func HandleConnection(conn netio.Ctx) {
	if err := handShake(conn); err != nil {
		logger.Info("Socks5 handshake error: %s", err.Error())
		return
	}
	addr, err := parseTarget(conn)
	if err != nil {
		logger.Info("socks consult transfer mode or parse target: %s", err.Error())
		return
	}
	pipeWhenClose(conn, addr)
}
