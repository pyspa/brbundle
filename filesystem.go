package brbundle

import (
	"os"
	"io"
	"net/http"
)

type FileEntry interface{
	RawReader() io.Reader
	BrotliReader() io.Reader
	Stat() os.FileInfo
}

type Pod interface {
	SetEncryptionKey(key []byte)
	FindEntry(path string) FileEntry
	ReadDirs(path string) []FileEntry

	Open(name string) (http.File, error)
}

type FileSystemOpt struct {
	EncryptionKey []byte
	EncryptionType
}
