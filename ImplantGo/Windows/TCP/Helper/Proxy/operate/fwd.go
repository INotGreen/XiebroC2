package operate

import (
	"main/Helper/Proxy/crypto"
	"main/Helper/Proxy/logger"
	"main/Helper/Proxy/netio"
	"main/Helper/Proxy/option"
	"net"
	"time"
)

func local2RemoteTCP(local string, remote string, lenc bool, renc bool) {
	listener, err := net.Listen("tcp", local)
	if err != nil {
		logger.Warn("Listen on %s error: %s", local, err.Error())
		return
	}
	defer listener.Close()

	for {
		logger.Info("Wait for connection on %s", local)

		localConn, err := listener.Accept()
		if err != nil {
			logger.Warn("Handle local connect error: %s", err.Error())
			continue
		}

		go func() {
			defer localConn.Close()

			logger.Info("Connection from %s", localConn.RemoteAddr().String())
			logger.Info("Connecting " + remote)

			localConnCtx, err := netio.NewTCPCtx(localConn, lenc)
			if err != nil {
				logger.Warn("Handle local connect error: %s", err.Error())
				return
			}

			remoteConn, err := net.DialTimeout(
				"tcp", remote,
				time.Millisecond*time.Duration(option.TIMEOUT),
			)
			if err != nil {
				logger.Warn("Connect remote %s error: %s", remote, err.Error())
				return
			}
			defer remoteConn.Close()

			remoteConnCtx, err := netio.NewTCPCtx(remoteConn, renc)
			if err != nil {
				logger.Warn("Connect remote %s error: %s", remote, err.Error())
				return
			}

			logger.Info("Open pipe: %s <== FWD ==> %s",
				localConn.RemoteAddr().String(), remoteConn.RemoteAddr().String())
			netio.PipeForward(localConnCtx, remoteConnCtx)
			logger.Info("Close pipe: %s <== FWD ==> %s",
				localConn.RemoteAddr().String(), remoteConn.RemoteAddr().String())
		}()
	}

}

func local2RemoteUDP(local string, remote string, lenc bool, renc bool) {
	localAddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", local, err.Error())
		return
	}
	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		logger.Warn("Listen udp on %s error: %s", local, err.Error())
		return
	}
	defer listener.Close()

	remoteAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", local, err.Error())
		return
	}
	remoteConn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		logger.Warn("Dial remote udp %s error: %s", local, err.Error())
		return
	}
	defer remoteConn.Close()

	listenerCtx, err := netio.NewUDPCtx(listener, lenc, false)
	if err != nil {
		return
	}
	remoteCtx, err := netio.NewUDPCtx(remoteConn, renc, true)
	if err != nil {
		return
	}

	netio.ForwardUDP(listenerCtx, remoteCtx)
}

func Local2Remote(local string, remote string, lenc bool, renc bool) {
	if option.PROTOCOL == "TCP" {
		logger.Success("Forward TCP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			local, lenc, remote, renc)
		local2RemoteTCP(local, remote, lenc, renc)
	} else {
		logger.Success("Forward UDP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			local, lenc, remote, renc)
		local2RemoteUDP(local, remote, lenc, renc)
	}
}

func local2LocalTCP(localA string, localB string, laenc bool, lbenc bool) {
	var listenerA net.Listener
	var listenerB net.Listener

	for {
		signal := make(chan byte)
		var localConnA net.Conn
		var localConnB net.Conn

		go func() {
			var err error
			listenerA, err = net.Listen("tcp", localA)
			if err != nil {
				logger.Warn("Listen on %s error: %s", localA, err.Error())
				return
			}
			defer listenerA.Close()

			for {
				logger.Info("Wait for connection on %s", localA)

				var err error
				localConnA, err = listenerA.Accept()
				if err != nil {
					logger.Warn("Handle connection error: %s", err.Error())
					continue
				}
				break
			}
			signal <- 'A'
		}()

		go func() {
			var err error
			listenerB, err = net.Listen("tcp", localB)
			if err != nil {
				logger.Warn("Listen on %s error: %s", localB, err.Error())
				return
			}
			defer listenerB.Close()

			for {
				logger.Info("Wait for connection on %s", localB)

				var err error
				localConnB, err = listenerB.Accept()
				if err != nil {
					logger.Warn("Handle connection error: %s", err.Error())
					continue
				}
				break
			}
			signal <- 'B'
		}()

		switch <-signal {
		case 'A':
			logger.Info("%s connected, waiting for %s", localA, localB)
		case 'B':
			logger.Info("%s connected, waiting for %s", localB, localA)
		}

		<-signal

		go func() {
			defer func() {
				if localConnA != nil {
					localConnA.Close()
				}

				if localConnB != nil {
					localConnB.Close()
				}
			}()

			localConnCtxA, err := netio.NewTCPCtx(localConnA, laenc)
			if err != nil {
				logger.Warn("handle local %s error: %s", localA, err.Error())
			}

			localConnCtxB, err := netio.NewTCPCtx(localConnB, lbenc)
			if err != nil {
				logger.Warn("handle local %s error: %s", localB, err.Error())
			}

			logger.Info("Open pipe: %s <== FWD ==> %s",
				localConnA.RemoteAddr().String(), localConnB.RemoteAddr().String())
			netio.PipeForward(localConnCtxA, localConnCtxB)
			logger.Info("Close pipe: %s <== FWD ==> %s",
				localConnA.RemoteAddr().String(), localConnB.RemoteAddr().String())
		}()
	}
}

func local2LocalUDP(localA string, localB string, laenc bool, lbenc bool) {
	localAddrA, err := net.ResolveUDPAddr("udp", localA)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", localA, err.Error())
		return
	}
	listenerA, err := net.ListenUDP("udp", localAddrA)
	if err != nil {
		logger.Warn("Listen udp on %s error: %s", localA, err.Error())
		return
	}
	defer listenerA.Close()

	localAddrB, err := net.ResolveUDPAddr("udp", localB)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", localB, err.Error())
		return
	}
	listenerB, err := net.ListenUDP("udp", localAddrB)
	if err != nil {
		logger.Warn("Listen udp on %s error: %s", localB, err.Error())
		return
	}
	defer listenerB.Close()

	listenerCtxA, err := netio.NewUDPCtx(listenerA, laenc, false)
	if err != nil {
		return
	}
	listenerCtxB, err := netio.NewUDPCtx(listenerB, lbenc, false)
	if err != nil {
		return
	}

	netio.ForwardUnconnectedUDP(listenerCtxA, listenerCtxB)
}

func Local2Local(localA string, localB string, laenc bool, lbenc bool) {
	if option.PROTOCOL == "TCP" {
		logger.Success("Forward TCP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			localA, laenc, localB, lbenc)

		local2LocalTCP(localA, localB, laenc, lbenc)
	} else {
		logger.Success("Forward UDP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			localA, laenc, localB, lbenc)
		local2LocalUDP(localA, localB, laenc, lbenc)
	}
}

func remote2remoteTCP(remoteA string, remoteB string, raenc bool, rbenc bool) {
	for {
		var remoteConnA net.Conn
		var remoteConnB net.Conn

		signal := make(chan struct{})

		go func() {
			for {
				var err error
				logger.Info("Connecting remote %s", remoteA)

				remoteConnA, err = net.DialTimeout(
					"tcp", remoteA,
					time.Millisecond*time.Duration(option.TIMEOUT),
				)
				if err != nil {
					logger.Info("Connect remote %s error, retrying", remoteA)
					time.Sleep(option.CONNECTING_RETRY_DURATION * time.Millisecond)
					continue
				}
				break
			}

			signal <- struct{}{}
		}()

		go func() {
			for {
				var err error
				logger.Info("Connecting remote %s", remoteB)

				remoteConnB, err = net.DialTimeout(
					"tcp", remoteB,
					time.Millisecond*time.Duration(option.TIMEOUT),
				)
				if err != nil {
					logger.Info("Connect remote %s error, retrying", remoteB)
					time.Sleep(option.CONNECTING_RETRY_DURATION * time.Millisecond)
					continue
				}
				break
			}

			signal <- struct{}{}
		}()

		<-signal
		<-signal

		go func() {
			defer func() {
				if remoteConnA != nil {
					remoteConnA.Close()
				}

				if remoteConnB != nil {
					remoteConnB.Close()
				}
			}()

			if remoteConnA != nil && remoteConnB != nil {
				remoteConnCtxA, err := netio.NewTCPCtx(remoteConnA, raenc)
				if err != nil {
					logger.Warn("Handle remote %s error: %s", remoteA, err.Error())
				}
				remoteConnCtxB, err := netio.NewTCPCtx(remoteConnB, rbenc)
				if err != nil {
					logger.Warn("Handle remote %s error: %s", remoteB, err.Error())
				}

				logger.Info("Start pipe: %s <== FWD ==> %s",
					remoteConnA.RemoteAddr().String(), remoteConnB.RemoteAddr().String())
				netio.PipeForward(remoteConnCtxA, remoteConnCtxB)
				logger.Info("Close pipe: %s <== FWD ==> %s",
					remoteConnA.RemoteAddr().String(), remoteConnB.RemoteAddr().String())
			}
		}()
	}
}

func remote2remoteUDP(remoteA string, remoteB string, raenc bool, rbenc bool) {
	remoteAddrA, err := net.ResolveUDPAddr("udp", remoteA)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", remoteA, err.Error())
		return
	}
	remoteConnA, err := net.DialUDP("udp", nil, remoteAddrA)
	if err != nil {
		logger.Warn("Dial remote udp %s error: %s", remoteA, err.Error())
		return
	}
	defer remoteConnA.Close()

	remoteAddrB, err := net.ResolveUDPAddr("udp", remoteB)
	if err != nil {
		logger.Warn("Parse udp address %s error: %s", remoteB, err.Error())
		return
	}
	remoteConnB, err := net.DialUDP("udp", nil, remoteAddrB)
	if err != nil {
		logger.Warn("Dial remote udp %s error: %s", remoteB, err.Error())
		return
	}
	defer remoteConnB.Close()

	remoteCtxA, err := netio.NewUDPCtx(remoteConnA, raenc, true)
	if err != nil {
		return
	}
	remoteCtxB, err := netio.NewUDPCtx(remoteConnB, rbenc, true)
	if err != nil {
		return
	}

	{
		// Need to send init packet to register the remote address, it doesn't matter even tough target is not `iox`
		//
		// There is a design fault here, and I need to consider the case where the FORWARD_WITHOUT_DEC flag is set
		// but actually needs to be encrypted, otherwise there is no IV in the ciphertext
		if raenc {
			iv, err := crypto.RandomNonce()
			cipher, err := crypto.NewCipher(iv)
			if err != nil {
				return
			}

			b := make([]byte, 4, 20)
			copy(b, netio.UDP_INIT_PACKET)

			cipher.StreamXOR(b, b)
			b = append(b, iv...)
			remoteCtxA.Write(b)

		} else {
			remoteCtxA.Write(netio.UDP_INIT_PACKET)
		}
		if rbenc {
			iv, err := crypto.RandomNonce()
			cipher, err := crypto.NewCipher(iv)
			if err != nil {
				return
			}

			b := make([]byte, 4, 20)
			copy(b, netio.UDP_INIT_PACKET)

			cipher.StreamXOR(b, b)
			b = append(b, iv...)
			remoteCtxB.Write(b)

		} else {
			remoteCtxB.Write(netio.UDP_INIT_PACKET)
		}
	}

	netio.ForwardUDP(remoteCtxA, remoteCtxB)
}

func Remote2Remote(remoteA string, remoteB string, raenc bool, rbenc bool) {
	if option.PROTOCOL == "TCP" {
		logger.Success("Forward TCP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			remoteA, raenc, remoteB, rbenc)
		remote2remoteTCP(remoteA, remoteB, raenc, rbenc)
	} else {
		logger.Success("Forward UDP traffic between %s (encrypted: %v) and %s (encrypted: %v)",
			remoteA, raenc, remoteB, rbenc)
		remote2remoteUDP(remoteA, remoteB, raenc, rbenc)
	}
}
