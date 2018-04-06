package brbundle

import (
	"errors"
)

type CompressionType int
type EncryptionType int

const (
	NoCompression CompressionType = iota
	Brotli
	LZ4

	NoEncryption EncryptionType = iota
	AES
	ChaCha20Poly1305
)

func (c CompressionType) String() string {
	switch c {
	case Brotli:
		return "brotli"
	case LZ4:
		return "lz4"
	case NoCompression:
		return "no"
	}
	return ""
}

func (c CompressionType) VariableName() string {
	switch c {
	case Brotli:
		return "brbundle.Brotli"
	case LZ4:
		return "brbundle.LZ4"
	case NoCompression:
		return "brbundle.NoCompression"
	}
	return ""
}


func (e EncryptionType) String() string {
	switch e {
	case AES:
		return "AES-256-GCM"
	case ChaCha20Poly1305:
		return "ChaCha20-Poly1305"
	case NoEncryption:
		return "no"
	}
	return ""
}

func (e EncryptionType) VariableName() string {
	switch e {
	case AES:
		return "brbundle.AES"
	case ChaCha20Poly1305:
		return "brbundle.ChaCha20Poly1305"
	case NoEncryption:
		return "brbundle.NoEncryption"
	}
	return ""
}


type FileSystem int

type PresetConfig []FileSystem

const (
	Embedded FileSystem = iota
	Append
	WorkingDir
	HomeDir
	Injected
)

var Production = []FileSystem{
	Embedded,
	Append,
}

var Development = []FileSystem{
	Injected,
	WorkingDir,
	Embedded,
	Append,
}

type Bundle struct {
	pods []Pod
}

func (b *Bundle) AddPod(pod Pod) {
	b.pods = append(b.pods, pod)
}

func (b Bundle) FindEntry(path string) (FileEntry, error) {
	for _, pod := range b.pods {
		entry := pod.FindEntry(path)
		if entry != nil {
			return entry, nil
		}
	}
	return nil, errors.New("Can't read the file")
}

func (b Bundle) ReadDirs(path string) ([]FileEntry, error) {
	for _, pod := range b.pods {
		entries := pod.ReadDirs(path)
		if entries != nil {
			return entries, nil
		}
	}
	return nil, errors.New("Can't read the dir")
}

type Entry struct {
	Path string
	FileMode
	Content []byte
	Info
}