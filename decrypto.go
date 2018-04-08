package brbundle

import (
	"io"
	"crypto/aes"
	"crypto/cipher"
	"bytes"
	"io/ioutil"
	"golang.org/x/crypto/chacha20poly1305"
	"errors"
)

type Decryptor interface {
	Decrypto(input io.Reader) (io.Reader, error)
	SetKey(key []byte) error
	NeedKey() bool
	HasKey() bool
}

type aesDecryptor struct {
	aead cipher.AEAD
}

func (a aesDecryptor) Decrypto(input io.Reader) (io.Reader, error) {
	if a.aead == nil {
		return nil, errors.New("Encryption Key is not set. Call SetKey() or set it via 1st param of Pod")
	}
	nonce := make([]byte, a.aead.NonceSize())
	io.ReadFull(input, nonce)
	cipherData, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	plainData, err := a.aead.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(plainData), nil
}

func (a aesDecryptor) NeedKey() bool {
	return true
}

func (a *aesDecryptor) SetKey(key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	a.aead = aesgcm
	return nil
}

func (a aesDecryptor) HasKey() bool {
	return a.aead != nil
}

func AESDecryptor(key ...[]byte) Decryptor {
	result := &aesDecryptor{}
	if len(key) > 0 {
		result.SetKey(key[0])
	}
	return result
}

type chaChaDecryptor struct {
	aead cipher.AEAD
}

func (c chaChaDecryptor) Decrypto(input io.Reader) (io.Reader, error) {
	nonce := make([]byte, c.aead.NonceSize())
	io.ReadFull(input, nonce)
	cipherData, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	plainData, err := c.aead.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(plainData), nil
}

func (c chaChaDecryptor) NeedKey() bool {
	return true
}

func (c *chaChaDecryptor) SetKey(key []byte) error {
	chacha, err := chacha20poly1305.New(key)
	if err != nil {
		return err
	}
	c.aead = chacha
	return nil
}

func (c chaChaDecryptor) HasKey() bool {
	return c.aead != nil
}

func ChaChaDecryptor(key ...[]byte) Decryptor {
	result := &chaChaDecryptor{}
	if len(key) > 0 {
		result.SetKey(key[0])
	}
	return result
}

type nullDecryptor struct {}

func (n nullDecryptor) Decrypto(input io.Reader) (io.Reader, error) {
	return input, nil
}

func (n *nullDecryptor) SetKey(key []byte) error {
	return errors.New("NullDecryptor can't accept key")
}

func (n nullDecryptor) NeedKey() bool {
	return false
}

func (n nullDecryptor) HasKey() bool {
	return true
}

func NullDecryptor() Decryptor {
	return &nullDecryptor{}
}
