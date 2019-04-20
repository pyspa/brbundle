package brbundle

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type folderBundle struct {
	baseBundle
	rootFolder string
}

func newFolderBundle(rootFolder string, encrypted bool, option Option) *folderBundle {
	decryptorType := NotToEncrypto
	if encrypted {
		decryptorType = UseAES
	}
	return &folderBundle{
		baseBundle: baseBundle{
			mountPoint:    option.MountPoint,
			name:          option.Name,
			decryptorType: decryptorType,
		},
		rootFolder: rootFolder,
	}
}

func (f folderBundle) find(searchPath string) (FileEntry, error) {
	filePath := filepath.Join(f.rootFolder, searchPath)
	s, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	decryptor, err := f.baseBundle.getDecryptor()
	if err != nil {
		return nil, err
	}
	return &folderFileEntry{
		filePath:    filePath,
		logicalPath: path.Clean("/" + path.Join(f.baseBundle.mountPoint, searchPath)),
		info:        s,
		decryptor:   decryptor,
	}, nil
}

func (f folderBundle) readdir(path string) []FileEntry {
	panic("implement me")
}

func (f folderBundle) close() {
	// do nothing
}

type folderFileEntry struct {
	filePath    string
	info        os.FileInfo
	logicalPath string
	decryptor   Decryptor
}

func (f folderFileEntry) Reader() (io.ReadCloser, error) {
	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	decryptoReader, err := f.decryptor.Decrypto(file)
	if err != nil {
		return nil, err
	}
	return NewReadCloser(decryptoReader, file), nil
}

func (f folderFileEntry) BrotliReader() (io.ReadCloser, error) {
	return nil, errors.New("Source data is not compressed by brotli")
}

func (f folderFileEntry) Stat() os.FileInfo {
	return f.info
}

func (f folderFileEntry) Name() string {
	return path.Base(f.filePath)
}

func (f folderFileEntry) Path() string {
	return f.logicalPath
}

func (f folderFileEntry) Etag() string {
	size := int(f.info.Size())
	return fmt.Sprintf("%x-%x", size, f.info.ModTime().Unix())
}
