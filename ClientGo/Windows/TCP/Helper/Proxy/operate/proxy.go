package operate

import (
	"main/Helper/Proxy/logger"
	"main/Helper/Proxy/netio"
	"main/Helper/Proxy/socks5"
	"net"
	"os"
	"os/signal"
)

func ProxyLocal(local string, encrypted bool) {
	listener, err := net.Listen("tcp", local)
	if err != nil {
		logger.Warn("Socks5 listen on %s error: %s", local, err.Error())
		return
	}

	logger.Success("Start socks5 server on %s (encrypted: %v)", local, encrypted)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Warn("Socks5 handle local connect error: %s", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			connCtx, err := netio.NewTCPCtx(conn, encrypted)
			if err != nil {
				return
			}

			socks5.HandleConnection(connCtx)
		}()
	}
}

func ProxyRemote(remote string, encrypted bool) {
	session, ctlStream, err := clientHandshake(remote)
	if err != nil {
		logger.Warn(err.Error())
		return
	}
	defer session.Close()

	logger.Success("Remote socks5 handshake ok (encrypted: %v)", encrypted)

	connectRequest := make(chan uint8, MAX_CONNECTION)
	defer close(connectRequest)
	endSignal := make(chan struct{})

	// handle ctrl+C
	{
		sigs := make(chan os.Signal)
		signal.Notify(sigs, os.Interrupt)
		go func() {
			<-sigs
			ctlStream.Write(marshal(Protocol{
				CMD: CTL_CLEANUP,
				N:   0,
			}))
			logger.Success("Recv Ctrl+C, exit now")
			os.Exit(0)
		}()
	}

	// handle ctl stream
	go func() {
		defer ctlStream.Close()

		for {
			pb, err := readUntilEnd(ctlStream)
			if err != nil {
				logger.Warn("Control connection has been closed, exit now")
				os.Exit(-1)
			}

			p := unmarshal(pb)
			switch p.CMD {
			case CTL_CONNECT_ME:
				connectRequest <- p.N
			case CTL_CLEANUP:
				endSignal <- struct{}{}
				return
			}
		}
	}()

	// handle CONNECT_ME request
	for {
		select {
		case <-endSignal:
			logger.Success("Recv exit signal from remote, exit now")
			return
		case n := <-connectRequest:
			for n > 0 {
				go func() {
					stream, err := session.OpenStream()
					if err != nil {
						logger.Info(err.Error())
						return
					}
					defer stream.Close()

					connCtx, err := netio.NewTCPCtx(stream, encrypted)
					if err != nil {
						return
					}

					socks5.HandleConnection(connCtx)
				}()
				n--
			}
		}
	}
}

func ProxyRemoteL2L(control string, local string, cenc bool, lenc bool) {
	masterListener, err := net.Listen("tcp", control)
	if err != nil {
		logger.Warn("Listen on %s error", control)
		return
	}
	defer masterListener.Close()

	logger.Info("Listen on %s for reverse socks5", control)

	localListener, err := net.Listen("tcp", local)
	if err != nil {
		logger.Warn("Listen on %s error", local)
		return
	}
	defer localListener.Close()

	session, ctlStream, err := serverHandshake(masterListener)
	if err != nil {
		logger.Warn(err.Error())
		return
	}
	defer session.Close()
	defer ctlStream.Close()

	logger.Success("Reverse socks5 server handshake ok from %s (encrypted: %v)", session.RemoteAddr().String(), cenc)
	logger.Success("Socks5 server is listening on %s (encrypted: %v)", local, lenc)

	// handle ctrl+C
	{
		sigs := make(chan os.Signal)
		signal.Notify(sigs, os.Interrupt)
		go func() {
			<-sigs
			ctlStream.Write(marshal(Protocol{
				CMD: CTL_CLEANUP,
				N:   0,
			}))
			logger.Success("Recv Ctrl+C, exit now")
			os.Exit(0)
		}()
	}

	localConnBuffer := make(chan net.Conn, MAX_CONNECTION)
	defer close(localConnBuffer)

	// handle ctl stream read
	go func() {
		for {
			pb, err := readUntilEnd(ctlStream)
			if err != nil {
				logger.Warn("Control connection has been closed, exit now")
				os.Exit(-1)
			}

			p := unmarshal(pb)
			switch p.CMD {
			case CTL_CLEANUP:
				logger.Success("Recv exit signal from remote, exit now")
				os.Exit(0)
			}
		}
	}()

	// handle local connection
	go func() {
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				continue
			}

			localConnBuffer <- localConn

			_, err = ctlStream.Write(marshal(Protocol{
				CMD: CTL_CONNECT_ME,
				N:   1,
			}))
			if err != nil {
				logger.Warn("Control connection has been closed, exit now")
				os.Exit(-1)
			}
		}
	}()

	for {
		remoteStream, err := session.AcceptStream()
		if err != nil {
			continue
		}

		localConn := <-localConnBuffer

		go func() {
			defer remoteStream.Close()
			defer localConn.Close()

			remoteConnCtx, err := netio.NewTCPCtx(remoteStream, cenc)
			if err != nil {
				return
			}

			localConnCtx, err := netio.NewTCPCtx(localConn, lenc)
			if err != nil {
				return
			}

			netio.PipeForward(remoteConnCtx, localConnCtx)
		}()
	}
}
