package brbundle

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

type bundle interface {
	find(path string) (FileEntry, error)
	readdir(path string) []FileEntry
	close()
	setDecryptionKey(key string) error
	getName() string
}

type baseBundle struct {
	mountPoint    string
	name          string
	decryptorType string
	decryptor     Decryptor
}

func (b *baseBundle) setDecryptionKey(key string) error {
	switch b.decryptorType {
	case UseAES:
		byteKey, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return err
		}
		if len(byteKey) != (32 + 12) {
			return errors.New("Decoded key length is wrong")
		}
		decryptor, err := newAESDecryptor(byteKey)
		if err != nil {
			return err
		}
		b.decryptor = decryptor
	case NotToEncrypto:
		return fmt.Errorf("bundle '%s' is not encrypted", b.name)
	}
	return fmt.Errorf("bundle '%s' uses unknown encryption type", b.name)
}

func (b baseBundle) getDecryptor() (Decryptor, error) {
	if b.decryptor != nil {
		return b.decryptor, nil
	}
	switch b.decryptorType {
	case UseAES:
		return nil, fmt.Errorf("bundle '%s' is encrypted by AES but no key passed", b.name)
	case NotToEncrypto:
		b.decryptor = newNullDecryptor()
	}
	return b.decryptor, nil
}

type ReadCloser struct {
	reader io.Reader
	closer io.Closer
}

func (rc ReadCloser) Read(buf []byte) (int, error) {
	return rc.reader.Read(buf)
}

func (rc ReadCloser) Close() error {
	if rc.closer == nil {
		return nil
	}
	return rc.closer.Close()
}

func NewReadCloser(reader io.Reader, closer io.Closer) io.ReadCloser {
	return &ReadCloser{
		reader: reader,
		closer: closer,
	}
}

type FileEntry interface {
	Reader() (io.ReadCloser, error)
	BrotliReader() (io.ReadCloser, error)
	Stat() os.FileInfo
	Name() string
	Path() string
	Etag() string
}
