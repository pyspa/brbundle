package brbundle

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"io/ioutil"
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
		return nil, errors.New("Encryption Key is not set. Call SetKey() or set it via 1st param of Bundle")
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

type nullDecryptor struct{}

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
