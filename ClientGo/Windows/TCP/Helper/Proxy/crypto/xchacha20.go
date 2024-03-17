package crypto

import (
	"crypto/rand"
	"main/Helper/Proxy/crypto/chacha20"
)

var (
	SECRET_KEY []byte
	NONCE      []byte
)

func shuffle(bs []byte) {
	for i := range bs {
		bs[i] ^= byte(i) ^ bs[(i+1)%len(bs)]*((bs[len(bs)-1-i]*bs[i])%255)
	}
}

func ExpandKey(key []byte) {
	SECRET_KEY = make([]byte, 0x20)
	NONCE = make([]byte, 0x18)

	if len(key) < 0x20 {
		var c byte = 0x20 - byte(len(key)&0x1F)

		for i := 0; i < int(c); i++ {
			key = append(key, c)
		}
	}

	copy(SECRET_KEY, key[:0x20])
	copy(NONCE, append(key[:0xC], key[len(key)-0xC:]...))

	for i := range SECRET_KEY {
		SECRET_KEY[i] = (SECRET_KEY[i] + byte(i)%255)
	}

	shuffle(SECRET_KEY)
	shuffle(NONCE)
}

type Cipher struct {
	c *chacha20.Cipher
}

func NewCipherPair() (*Cipher, *Cipher, error) {
	ccA, err := chacha20.New(SECRET_KEY, NONCE)
	if err != nil {
		return nil, nil, err
	}
	ccB, err := chacha20.New(SECRET_KEY, NONCE)
	if err != nil {
		return nil, nil, err
	}

	return &Cipher{c: ccA}, &Cipher{c: ccB}, nil
}

func RandomNonce() ([]byte, error) {
	iv := make([]byte, 0x18)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

func NewCipher(nonce []byte) (*Cipher, error) {
	cc, err := chacha20.New(SECRET_KEY, nonce)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		c: cc,
	}, nil
}

func (c Cipher) StreamXOR(dst []byte, src []byte) {
	c.c.XORKeyStream(dst, src)
}
