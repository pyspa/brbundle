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

func (f folderBundle) close() {
	// do nothing
}

func (f folderBundle) dirs() []string {
	var dirs []string
	filepath.Walk(f.rootFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			rel, err := filepath.Rel(f.rootFolder, path)
			if err != nil {
				return err
			}
			if rel == "." {
				dirs = append(dirs, strings.TrimSuffix(f.mountPoint, "/"))
			} else if f.mountPoint == "" {
				dirs = append(dirs, strings.ReplaceAll(rel, `\`, "/"))
			} else {
				dirs = append(dirs, f.mountPoint+strings.ReplaceAll(rel, `\`, "/"))
			}
		}
		return nil
	})
	return dirs
}

func (f folderBundle) filesInDir(dirName string) []string {
	if !strings.HasPrefix(dirName, f.mountPoint) {
		return nil
	}
	dirName = dirName[len(f.mountPoint):]
	dirs, err := ioutil.ReadDir(filepath.Join(f.rootFolder, dirName))
	if err != nil {
		return nil
	}
	var result []string
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		result = append(result, dir.Name())
	}
	return result
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
		if m, err := mimetype.DetectFile(f.filePath); err != nil {
			f.contentType = m.String()
		} else {
			f.contentType = "application/octet-stream"
		}
	}
	return fmt.Sprintf("%x-%x", size, f.info.ModTime().Unix()), f.contentType
}

func (f folderFileEntry) GetLocalPath() (string, error) {
	return f.filePath, nil
}
