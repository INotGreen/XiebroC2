package Proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/Acebond/ReverseSocks5/statute"
)

// A Request represents request received by a server
type Request struct {
	statute.Request
	// AuthContext provided during negotiation
	// AuthContext *AuthContext
	// LocalAddr of the network server listen
	LocalAddr net.Addr
	// RemoteAddr of the network that sent the request
	RemoteAddr net.Addr
	// DestAddr of the actual destination (might be affected by rewrite)
	DestAddr *statute.AddrSpec
	// Reader connect of request
	Reader io.Reader
	// RawDestAddr of the desired destination
	RawDestAddr *statute.AddrSpec
}

// ParseRequest creates a new Request from the tcp connection
func ParseRequest(bufConn io.Reader) (*Request, error) {
	hd, err := statute.ParseRequest(bufConn)
	if err != nil {
		return nil, err
	}
	return &Request{
		Request:     hd,
		RawDestAddr: &hd.DstAddr,
		Reader:      bufConn,
	}, nil
}

// handleRequest is used for request processing after authentication
func handleRequest(write io.Writer, req *Request) error {

	// Resolve the address if we have a FQDN
	dest := req.RawDestAddr
	if dest.FQDN != "" {
		addr, err := net.ResolveIPAddr("ip", dest.FQDN)
		if err != nil {
			if err := SendReply(write, statute.RepHostUnreachable, nil); err != nil {
				return fmt.Errorf("failed to send reply, %v", err)
			}
			return fmt.Errorf("failed to resolve destination[%v], %v", dest.FQDN, err)
		}
		dest.IP = addr.IP
	}

	// Apply any address rewrites
	req.DestAddr = req.RawDestAddr

	// Switch on the command
	switch req.Command {

	case statute.CommandConnect:
		return handleConnect(write, req)

	case statute.CommandAssociate:
		return handleAssociate(write, req)

	default:
		if err := SendReply(write, statute.RepCommandNotSupported, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("unsupported command[%v]", req.Command)
	}
}

// handleConnect is used to handle a connect command
func handleConnect(writer io.Writer, request *Request) error {

	target, err := net.Dial("tcp", request.DestAddr.String())
	if err != nil {
		msg := err.Error()
		resp := statute.RepHostUnreachable
		if strings.Contains(msg, "refused") {
			resp = statute.RepConnectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = statute.RepNetworkUnreachable
		}
		if err := SendReply(writer, resp, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("connect to %v failed, %v", request.RawDestAddr, err)
	}
	defer target.Close()

	// Send success
	if err := SendReply(writer, statute.RepSuccess, target.LocalAddr()); err != nil {
		return fmt.Errorf("failed to send reply, %v", err)
	}

	// Start proxying
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		io.CopyBuffer(target, request.Reader, buf[:cap(buf)])
		wg.Done()
	}()
	go func() {
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		io.CopyBuffer(writer, target, buf[:cap(buf)])
		wg.Done()
	}()

	wg.Wait()

	return nil
}

// handleAssociate is used to handle a connect command
func handleAssociate(writer io.Writer, request *Request) error {

	bindLn, err := net.ListenUDP("udp", nil)
	if err != nil {
		if err := SendReply(writer, statute.RepServerFailure, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("listen udp failed, %v", err)
	}
	log.Printf("client want to used addr %v, listen addr: %s", request.DestAddr, bindLn.LocalAddr())
	// send BND.ADDR and BND.PORT, client used
	if err = SendReply(writer, statute.RepSuccess, bindLn.LocalAddr()); err != nil {
		return fmt.Errorf("failed to send reply, %v", err)
	}

	go func() {
		// read from client and write to remote server
		conns := sync.Map{}
		bufPool := bufferPool.Get()
		defer func() {
			bufferPool.Put(bufPool)
			bindLn.Close()
			conns.Range(func(key, value any) bool {
				if connTarget, ok := value.(net.Conn); !ok {
					log.Printf("conns has illegal item %v:%v", key, value)
				} else {
					connTarget.Close()
				}
				return true
			})
		}()
		for {
			n, srcAddr, err := bindLn.ReadFromUDP(bufPool[:cap(bufPool)])
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}
				continue
			}
			pk, err := statute.ParseDatagram(bufPool[:n])
			if err != nil {
				continue
			}

			// check src addr whether equal requst.DestAddr
			srcEqual := ((request.DestAddr.IP.IsUnspecified()) || request.DestAddr.IP.Equal(srcAddr.IP)) && (request.DestAddr.Port == 0 || request.DestAddr.Port == srcAddr.Port) //nolint:lll
			if !srcEqual {
				continue
			}

			connKey := srcAddr.String() + "--" + pk.DstAddr.String()

			if target, ok := conns.Load(connKey); !ok {
				// if the 'connection' doesn't exist, create one and store it
				targetNew, err := net.Dial("udp", pk.DstAddr.String())
				if err != nil {
					log.Printf("connect to %v failed, %v", pk.DstAddr, err)
					// TODO:continue or return Error?
					continue
				}
				conns.Store(connKey, targetNew)
				// read from remote server and write to original client
				go func() {
					bufPool := bufferPool.Get()
					defer func() {
						targetNew.Close()
						conns.Delete(connKey)
						bufferPool.Put(bufPool)
					}()

					for {
						buf := bufPool[:cap(bufPool)]
						n, err := targetNew.Read(buf)
						if err != nil {
							if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
								return
							}
							log.Printf("read data from remote %s failed, %v", targetNew.RemoteAddr().String(), err)
							return
						}
						tmpBufPool := bufferPool.Get()
						proBuf := tmpBufPool
						proBuf = append(proBuf, pk.Header()...)
						proBuf = append(proBuf, buf[:n]...)
						if _, err := bindLn.WriteTo(proBuf, srcAddr); err != nil {
							bufferPool.Put(tmpBufPool)
							log.Printf("write data to client %s failed, %v", srcAddr, err)
							return
						}
						bufferPool.Put(tmpBufPool)
					}
				}()
				if _, err := targetNew.Write(pk.Data); err != nil {
					log.Printf("write data to remote server %s failed, %v", targetNew.RemoteAddr().String(), err)
					return
				}
			} else {
				if _, err := target.(net.Conn).Write(pk.Data); err != nil {
					log.Printf("write data to remote server %s failed, %v", target.(net.Conn).RemoteAddr().String(), err)
					return
				}
			}
		}
	}()

	buf := bufferPool.Get()
	defer bufferPool.Put(buf)

	for {
		_, err := request.Reader.Read(buf[:cap(buf)])
		// sf.logger.Errorf("read data from client %s, %d bytesm, err is %+v", request.RemoteAddr.String(), num, err)
		if err != nil {
			bindLn.Close()
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
	}
}

// SendReply is used to send a reply message
// rep: reply status see statute's statute file
func SendReply(w io.Writer, rep uint8, bindAddr net.Addr) error {
	rsp := statute.Reply{
		Version:  statute.VersionSocks5,
		Response: rep,
		BndAddr: statute.AddrSpec{
			AddrType: statute.ATYPIPv4,
			IP:       net.IPv4zero,
			Port:     0,
		},
	}

	if rsp.Response == statute.RepSuccess {
		if tcpAddr, ok := bindAddr.(*net.TCPAddr); ok && tcpAddr != nil {
			rsp.BndAddr.IP = tcpAddr.IP
			rsp.BndAddr.Port = tcpAddr.Port
		} else if udpAddr, ok := bindAddr.(*net.UDPAddr); ok && udpAddr != nil {
			rsp.BndAddr.IP = udpAddr.IP
			rsp.BndAddr.Port = udpAddr.Port
		} else {
			rsp.Response = statute.RepAddrTypeNotSupported
		}

		if rsp.BndAddr.IP.To4() != nil {
			rsp.BndAddr.AddrType = statute.ATYPIPv4
		} else if rsp.BndAddr.IP.To16() != nil {
			rsp.BndAddr.AddrType = statute.ATYPIPv6
		}
	}
	// Send the message
	_, err := w.Write(rsp.Bytes())
	return err
}
