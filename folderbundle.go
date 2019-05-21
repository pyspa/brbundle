package brbundle

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	mountPoint := option.MountPoint
	if mountPoint != "" && !strings.HasSuffix(mountPoint, "/") {
		mountPoint = mountPoint + "/"
	}
	return &folderBundle{
		baseBundle: baseBundle{
			mountPoint:    mountPoint,
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
	contentType string
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

func (f folderFileEntry) ReadAll() ([]byte, error) {
	reader, err := f.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func (f folderFileEntry) BrotliReader() (io.ReadCloser, error) {
	return nil, errors.New("Source data is not compressed by brotli")
}

func (f folderFileEntry) CompressedSize() int64 {
	return -1
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

func (f folderFileEntry) EtagAndContentType() (string, string) {
	size := int(f.info.Size())
	if f.contentType == "" {
		f.contentType, _, _ = mimetype.DetectFile(f.filePath)
	}
	return fmt.Sprintf("%x-%x", size, f.info.ModTime().Unix()), f.contentType
}

func (f folderFileEntry) GetLocalPath() (string, error) {
	return f.filePath, nil
}
