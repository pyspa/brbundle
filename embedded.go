package brbundle

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"time"
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
	entry        *Entry
	decompressor Decompressor
	decryptor    Decryptor
}

func (e embeddedFileEntry) Reader() (io.ReadCloser, error) {
	decryptoReader, err := e.decryptor.Decrypto(bytes.NewReader(e.entry.Data))
	if err != nil {
		return nil, err
	}
	return NewReadCloser(e.decompressor.Decompress(decryptoReader), nil), nil
}

func (e embeddedFileEntry) BrotliReader() (io.ReadCloser, error) {
	if _, ok := e.decompressor.(*brotliDecompressor); ok {
		decryptoReader, err := e.decryptor.Decrypto(bytes.NewReader(e.entry.Data))
		if err != nil {
			return nil, err
		}
		return NewReadCloser(decryptoReader, nil), nil
	}
	return nil, errors.New("Source data is not compressed by brotli")
}

func (e embeddedFileEntry) Stat() os.FileInfo {
	return &embeddedFileInfo{e.entry}
}

func (e embeddedFileEntry) Name() string {
	return path.Base(e.entry.Path)
}

func (e embeddedFileEntry) Path() string {
	return e.entry.Path
}

type embeddedFileInfo struct {
	entry *Entry
}

func (e embeddedFileInfo) Name() string {
	return path.Base(e.entry.Path)
}

func (e embeddedFileInfo) Size() int64 {
	return e.entry.OriginalSize
}

func (e embeddedFileInfo) Mode() os.FileMode {
	return os.FileMode(e.entry.FileMode)
}
func (e embeddedFileInfo) ModTime() time.Time {
	return e.entry.Mtime
}

func (e embeddedFileInfo) IsDir() bool {
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

func (e EmbeddedPod) Find(path string) FileEntry {
	entry, ok := e.files[path]
	if ok {
		return &embeddedFileEntry{
			entry:        entry,
			decompressor: e.decompressor,
			decryptor:    e.decryptor,
		}
	}
	return nil
}

func (e EmbeddedPod) Readdir(path string) []FileEntry {
	filePaths := e.dirs[path]
	var result []FileEntry
	for _, filePath := range filePaths {
		entry := e.Find(filePath)
		if entry != nil {
			result = append(result, entry)
		}
	}
	return result
}

func (e EmbeddedPod) Close() error {
	// do nothing
	return nil
}

func NewEmbeddedPod(decompressor Decompressor, decryptor Decryptor, dirs map[string][]string, files map[string]*Entry) func(key ...[]byte) (FilePod, error) {
	return func(key ...[]byte) (FilePod, error) {
		if decryptor.NeedKey() {
			if len(key) < 1 {
				return nil, errors.New("Key to decrypto is needed")
			}
			decryptor.SetKey(key[0])
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
		if decryptor.NeedKey() {
			if len(key) < 1 {
				panic(errors.New("Key to decrypto is needed"))
			}
			decryptor.SetKey(key[0])
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
		return pod
	}
}
