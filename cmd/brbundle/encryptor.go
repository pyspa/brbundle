package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"

	"github.com/shibukawa/brbundle"
	"golang.org/x/crypto/chacha20poly1305"
)

type Encryptor struct {
	etype  brbundle.EncryptionType
	aead   cipher.AEAD
	key    []byte
	nonce  []byte
	reader *io.PipeReader
	writer *io.PipeWriter
	size   int

	processingPath string
}

func NewEncryptor(etype brbundle.EncryptionType, key, nonce []byte) *Encryptor {
	e := &Encryptor{
		etype,
		nil,
		key,
		nonce,
		nil,
		nil,
		0,
		"",
	}
	switch e.etype {
	case brbundle.AES:
		block, err := aes.NewCipher(e.key)
		if err != nil {
			panic(err.Error())
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}
		e.aead = gcm
	case brbundle.ChaCha20Poly1305:
		chacha, err := chacha20poly1305.New(e.key)
		if err != nil {
			panic(err.Error())
		}
		e.aead = chacha
	}
	return e
}

func (e *Encryptor) Init() {
	reader, writer := io.Pipe()
	e.reader = reader
	e.writer = writer
	e.size = 0
}

func (e Encryptor) Size() int {
	return e.size
}

func (e *Encryptor) SetPath(path string) {
	e.processingPath = path
}

func (e *Encryptor) Write(data []byte) (n int, err error) {
	n, err = e.writer.Write(data)
	e.size = e.size + n
	return
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
		var nonce []byte
		if len(e.nonce) == 0 {
			nonce = make([]byte, e.aead.NonceSize())
			_, err = rand.Read(nonce)
			if err != nil {
				panic("nonce generation error")
			}
		} else {
			nonce = e.nonce
		}
		_, err = w.Write(nonce)
		if err != nil {
			return 0, err
		}
		cipherContent := e.aead.Seal(nil, nonce, src, nil)
		n, err := w.Write(cipherContent)
		return int64(e.aead.NonceSize() + n), err
	}
	return
}
