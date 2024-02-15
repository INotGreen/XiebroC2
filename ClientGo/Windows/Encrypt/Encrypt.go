package Encrypt

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"io/ioutil"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

var AesKey = "QWERt_CSDMAHUATE"

func Compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}

func Encrypt(data []byte) ([]byte, error) {
	return aesECBEncrypt(data, []byte(AesKey))
}

func Decrypt(data []byte) ([]byte, error) {
	return aesECBDncrypt(data, []byte(AesKey))
}

// AES加密
func aesECBEncrypt(pt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return pt, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
	if err != nil {
		return pt, err
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

// AES解密
func aesECBDncrypt(ct, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ct, err
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		return ct, err
	}
	return pt, nil
}
