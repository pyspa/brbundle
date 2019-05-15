package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"

	"go.pyspa.org/brbundle"
)

type Encryptor struct {
	etype  brbundle.EncryptionType
	aead   cipher.AEAD
	nonce  []byte
	reader *io.PipeReader
	writer *io.PipeWriter

	processingPath string
}

func decodeEncryptKey(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}
	bytesKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, errors.New("Decode base64 error. Use brbundle generate-key command to generate key.")
	}
	if len(bytesKey) != (32 + 12) {
		return nil, errors.New("Encryption-key length is wrong. Use brbundle generate-key command to generate key.")
	}
	return bytesKey, nil
}

func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) == 0 {
		return &Encryptor{
			etype: brbundle.NoEncryption,
		}, nil
	} else {
		block, err := aes.NewCipher(key[:32])
		if err != nil {
			panic(err.Error())
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}
		return &Encryptor{
			etype: brbundle.AES,
			aead:  gcm,
			nonce: key[32:],
		}, nil
	}
}

func (e Encryptor) EncryptionFlag() string {
	if e.etype == brbundle.NoEncryption {
		return brbundle.NotToEncrypto
	} else if e.etype == brbundle.AES {
		return brbundle.UseAES
	}
	panic("undefined encryption flags")
}

func (e *Encryptor) Init() {
	reader, writer := io.Pipe()
	e.reader = reader
	e.writer = writer
}

func (e *Encryptor) SetPath(path string) {
	e.processingPath = path
}

func (e *Encryptor) Write(data []byte) (n int, err error) {
	return e.writer.Write(data)
}

func (e *Encryptor) Close() {
	e.writer.Close()
}

func (e *Encryptor) WriteTo(w io.Writer) (n int64, err error) {
	if e.etype == brbundle.NoEncryption {
		n, err = io.Copy(w, e.reader)
	} else {
		src, err := ioutil.ReadAll(e.reader)
		if err != nil {
			return 0, err
		}
		cipherContent := e.aead.Seal(nil, e.nonce, src, nil)
		n, err := w.Write(cipherContent)
		return int64(n), err
	}
	return
}
