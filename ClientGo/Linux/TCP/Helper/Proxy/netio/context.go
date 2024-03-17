package netio

import (
	"main/Helper/Proxy/crypto"
	"main/Helper/Proxy/option"
	"net"
)

type Ctx interface {
	DecryptRead(b []byte) (int, error)
	EncryptWrite(b []byte) (int, error)

	net.Conn
}

var _ Ctx = &TCPCtx{}
var _ Ctx = &UDPCtx{}

type TCPCtx struct {
	net.Conn
	encrypted bool

	// Ensure stream cipher synchronous
	encCipher *crypto.Cipher
	decCipher *crypto.Cipher
}

func NewTCPCtx(conn net.Conn, encrypted bool) (*TCPCtx, error) {
	// if tc, ok := conn.(*net.TCPConn); ok {
	//     tc.SetLinger(0)
	// }

	encrypted = encrypted && !option.FORWARD_WITHOUT_DEC

	ctx := &TCPCtx{
		Conn:      conn,
		encrypted: encrypted,
	}

	if encrypted {
		encCipher, decCipher, err := crypto.NewCipherPair()
		if err != nil {
			return nil, err
		}

		ctx.encCipher = encCipher
		ctx.decCipher = decCipher
	}

	return ctx, nil
}

func (c *TCPCtx) DecryptRead(b []byte) (int, error) {
	n, err := c.Read(b)
	if err != nil {
		return n, err
	}

	if c.encrypted {
		c.decCipher.StreamXOR(b[:n], b[:n])
	}

	return n, err
}

func (c *TCPCtx) EncryptWrite(b []byte) (int, error) {
	if c.encrypted {
		c.encCipher.StreamXOR(b, b)
	}
	return c.Write(b)
}

type UDPCtx struct {
	*net.UDPConn
	encrypted  bool
	connected  bool
	remoteAddr *net.UDPAddr

	// sync.Mutex
}

func NewUDPCtx(conn *net.UDPConn, encrypted bool, connected bool) (*UDPCtx, error) {
	encrypted = encrypted && !option.FORWARD_WITHOUT_DEC

	ctx := &UDPCtx{
		UDPConn:   conn,
		encrypted: encrypted,
		connected: connected,
	}

	return ctx, nil
}

// Encryption for packet is different from stream
func (c *UDPCtx) DecryptRead(b []byte) (int, error) {
	var n int
	var err error

	if !c.connected {
		var remoteAddr *net.UDPAddr
		n, remoteAddr, err = c.ReadFromUDP(b)
		if err != nil {
			return n, err
		}
		c.remoteAddr = remoteAddr

	} else {
		n, err = c.Read(b)
		if err != nil {
			return n, err
		}
	}

	if c.encrypted {
		if len(b) < 0x18 {
			// no nonce, skip
			return 0, nil
		}
		nonce := b[n-0x18 : n]
		b = b[:n-0x18]

		cipher, err := crypto.NewCipher(nonce)
		if err != nil {
			return 0, err
		}

		n -= 0x18
		cipher.StreamXOR(b[:n], b[:n])
	}

	return n, err
}

func (c *UDPCtx) EncryptWrite(b []byte) (int, error) {
	if c.encrypted {
		iv, err := crypto.RandomNonce()
		cipher, err := crypto.NewCipher(iv)
		if err != nil {
			return 0, err
		}

		cipher.StreamXOR(b, b)
		b = append(b, iv...)
	}

	if !c.connected {
		return c.WriteTo(b, c.remoteAddr)
	}
	return c.Write(b)
}

/*
func (c UDPCtx) IsRemoteAddrRegistered() bool {
	return c.remoteAddr != nil
}
*/
