package brbundle

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type ZipBundle struct {
	decompressor Decompressor
	decryptor    Decryptor
	dirs         map[string][]string
	files        map[string]*zip.File
	onClose      func() error
}

func (z ZipBundle) Find(path string) FileEntry {
	entry, ok := z.files[path]
	if ok {
		return &zipFileEntry{
			entry:        entry,
			decompressor: z.decompressor,
			decryptor:    z.decryptor,
		}
	}
	return nil
}

func (z ZipBundle) Readdir(path string) []FileEntry {
	return nil
}

func (z *ZipBundle) OnClose(onClose func() error) {
	z.onClose = onClose
}

func (z ZipBundle) Close() error {
	if z.onClose == nil {
		return nil
	}
	return z.onClose()
}

type zipFileEntry struct {
	entry        *zip.File
	decompressor Decompressor
	decryptor    Decryptor
}

func (z zipFileEntry) Reader() (io.ReadCloser, error) {
	reader, err := z.entry.Open()
	if err != nil {
		return nil, err
	}
	decryptoReader, err := z.decryptor.Decrypto(reader)
	if err != nil {
		return nil, err
	}
	return NewReadCloser(z.decompressor.Decompress(decryptoReader), reader), nil
}

func (z zipFileEntry) BrotliReader() (io.ReadCloser, error) {
	if _, ok := z.decompressor.(*brotliDecompressor); ok {
		reader, err := z.entry.Open()
		if err != nil {
			return nil, err
		}
		decryptoReader, err := z.decryptor.Decrypto(reader)
		if err != nil {
			return nil, err
		}
		return NewReadCloser(decryptoReader, reader), nil
	}
	return nil, errors.New("Source data is not compressed by brotli")
}

func (z *zipFileEntry) Stat() os.FileInfo {
	originalSize, _ := strconv.ParseInt(strings.Split(z.entry.Comment, "-")[0], 16, 64)
	return &zipFileInfo{z.Stat(), originalSize}
}

func (z *zipFileEntry) Name() string {
	return path.Base(z.entry.Name)
}

func (z *zipFileEntry) Path() string {
	return "/" + z.entry.Name
}

type zipFileInfo struct {
	info         os.FileInfo
	originalSize int64
}

func (z *zipFileInfo) Name() string {
	return z.info.Name()
}

func (z *zipFileInfo) Size() int64 {
	return z.originalSize
}

func (z *zipFileInfo) Mode() os.FileMode {
	return z.info.Mode()
}
func (z *zipFileInfo) ModTime() time.Time {
	return z.info.ModTime()
}

func (z *zipFileInfo) IsDir() bool {
	return z.info.IsDir()
}

func (z *zipFileInfo) Sys() interface{} {
	return z.info.Sys()
}

func NewZipBundle(decompressor Decompressor, decryptor Decryptor, zipFilePath string) (FileBundle, error) {
	file, err := os.Open(zipFilePath)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	bundle, err := NewZipBundleFromReader(decompressor, decryptor, file, info.Size())
	if err != nil {
		bundle.(*ZipBundle).OnClose(func() error {
			return file.Close()
		})
	}
	return bundle, err
}

func MustZipBundle(decompressor Decompressor, decryptor Decryptor, zipFilePath string) FileBundle {
	bundle, err := NewZipBundle(decompressor, decryptor, zipFilePath)
	if err != nil {
		panic(err)
	}
	return bundle
}

func NewZipBundleFromReader(decompressor Decompressor, decryptor Decryptor, fileReader io.ReaderAt, size int64) (FileBundle, error) {
	if decryptor.NeedKey() && !decryptor.HasKey() {
		return nil, errors.New("Key to decrypto is needed")
	}
	reader, err := zip.NewReader(fileReader, size)
	if err != nil {
		return nil, err
	}
	dirs := make(map[string][]string)
	files := make(map[string]*zip.File)
	for _, file := range reader.File {
		files["/"+file.Name] = file
	}
	bundle := &ZipBundle{
		decompressor: decompressor,
		decryptor:    decryptor,
		dirs:         dirs,
		files:        files,
	}
	return bundle, nil
}

func MustZipBundleFromReader(decompressor Decompressor, decryptor Decryptor, reader io.ReaderAt, size int64) FileBundle {
	bundle, err := NewZipBundleFromReader(decompressor, decryptor, reader, size)
	if err != nil {
		panic(err)
	}
	return bundle
}

func NewZipBundleFromZipReader(decompressor Decompressor, decryptor Decryptor, reader *zip.Reader) (FileBundle, error) {
	if decryptor.NeedKey() && !decryptor.HasKey() {
		return nil, errors.New("Key to decrypto is needed")
	}
	dirs := make(map[string][]string)
	files := make(map[string]*zip.File)
	for _, file := range reader.File {
		files["/"+file.Name] = file
	}
	bundle := &ZipBundle{
		decompressor: decompressor,
		decryptor:    decryptor,
		dirs:         dirs,
		files:        files,
	}
	return bundle, nil
}

func MustZipBundleFromZipReader(decompressor Decompressor, decryptor Decryptor, reader *zip.Reader) FileBundle {
	bundle, err := NewZipBundleFromZipReader(decompressor, decryptor, reader)
	if err != nil {
		panic(err)
	}
	return bundle
}
