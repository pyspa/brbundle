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
}

type aesDecryptor struct {
	aead  cipher.AEAD
	nonce []byte
}

func (a aesDecryptor) Decrypto(input io.Reader) (io.Reader, error) {
	if a.aead == nil {
		return nil, errors.New("Encryption Key is not set. Call SetKey() or set it via 1st param of Bundle")
	}
	cipherData, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	plainData, err := a.aead.Open(nil, a.nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(plainData), nil
}

func newAESDecryptor(key []byte) (Decryptor, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	result := &aesDecryptor{
		aead:  aesgcm,
		nonce: key[32:],
	}
	return result, nil
}

type nullDecryptor struct{}

func (n nullDecryptor) Decrypto(input io.Reader) (io.Reader, error) {
	return input, nil
}

func newNullDecryptor() Decryptor {
	return &nullDecryptor{}
}
