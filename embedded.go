package brbundle

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
	"errors"
)

type Entry struct {
	Path         string
	FileMode     uint32
	OriginalSize int64
	Mtime        time.Time
	Data         []byte
	ETag         string
}

type embeddedFileEntry struct {
	entry *Entry
	decompressor  Decompressor
	decryptor     Decryptor
}

func (e embeddedFileEntry) RawReader() (io.Reader, error) {
	decryptoReader, err :=  e.decryptor.Decrypto(bytes.NewReader(e.entry.Data))
	if err != nil {
		return nil, err
	}
	return e.decompressor.Decompress(decryptoReader), nil
}

func (e embeddedFileEntry) BrotliReader() (io.Reader, error) {
	if _, ok := e.decompressor.(*brotliDecompressor); ok {
		decryptoReader, err :=  e.decryptor.Decrypto(bytes.NewReader(e.entry.Data))
		if err != nil {
			return nil, err
		}
		return decryptoReader, nil
	}
	return nil, errors.New("Source data is not compressed by brotli")
}

func (e *embeddedFileEntry) Stat() os.FileInfo {
	return &embeddedFileInfo{e.entry}
}

func (e *embeddedFileEntry) Name() string {
	return path.Base(e.entry.Path)
}

type embeddedFileInfo struct {
	entry *Entry
}

func (e *embeddedFileInfo) Name() string {
	return path.Base(e.entry.Path)
}

func (e *embeddedFileInfo) Size() int64 {
	return e.entry.OriginalSize
}

func (e *embeddedFileInfo) Mode() os.FileMode {
	return os.FileMode(e.entry.FileMode)
}
func (e *embeddedFileInfo) ModTime() time.Time {
	return e.entry.Mtime
}

func (e *embeddedFileInfo) IsDir() bool {
	return false
}

func (e *embeddedFileInfo) Sys() interface{} {
	return nil
}

type EmbeddedPod struct {
	decompressor  Decompressor
	decryptor     Decryptor
	dirs          map[string][]string
	files         map[string]*Entry
	encryptionKey []byte
}

func (e *EmbeddedPod) Find(path string) FileEntry {
	entry, ok := e.files[path]
	if ok {
		return &embeddedFileEntry{
			entry: entry,
			decompressor: e.decompressor,
			decryptor: e.decryptor,
		}
	}
	return nil
}

func (e *EmbeddedPod) Readdir(path string) []FileEntry {
	return nil
}

func (e *EmbeddedPod) Open(name string) (http.File, error) {
	return nil, nil
}

func NewEmbeddedPod(decompressor Decompressor, decryptor Decryptor, dirs map[string][]string, files map[string]*Entry) func(key ...[]byte) (FilePod, error) {
	return func(key ...[]byte) (FilePod, error) {
		if decryptor.NeedKey() && len(key) < 1 {
			return nil, fmt.Errorf("Key to decrypto is needed")
		}
		pod := &EmbeddedPod{
			decompressor: decompressor,
			decryptor:    decryptor,
			dirs:         dirs,
			files:        files,
		}
		if len(key) > 0 {
			pod.encryptionKey = key[0]
		}
		return pod, nil
	}
}

func MustEmbeddedPod(decompressor Decompressor, decryptor Decryptor, dirs map[string][]string, files map[string]*Entry) func(key ...[]byte) FilePod {
	return func(key ...[]byte) FilePod {
		if decryptor.NeedKey() && len(key) < 1 {
			panic(fmt.Errorf("Key to decrypto is needed"))
		}
		pod := &EmbeddedPod{
			decompressor: decompressor,
			decryptor:  decryptor,
			dirs:            dirs,
			files:           files,
		}
		if len(key) > 0 {
			pod.encryptionKey = key[0]
		}
		return pod
	}
}
