package brbundle

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type packedBundle struct {
	baseBundle
	reader *zip.Reader
	closer io.Closer
}

func newPackedBundle(reader *zip.Reader, closer io.Closer, option Option) *packedBundle {
	sort.Slice(reader.File, func(i, j int) bool {
		return reader.File[i].Name < reader.File[j].Name
	})
	result := &packedBundle{
		baseBundle: baseBundle{
			mountPoint:    option.MountPoint,
			name:          option.Name,
			decryptorType: reader.Comment[:1],
		},
		reader: reader,
		closer: closer,
	}
	return result
}

func (p packedBundle) find(searchPath string) (FileEntry, error) {
	i := sort.Search(len(p.reader.File), func(i int) bool {
		return p.reader.File[i].Name >= searchPath
	})
	if i < len(p.reader.File) && p.reader.File[i].Name == searchPath {
		decryptor, err := p.baseBundle.getDecryptor()
		if err != nil {
			return nil, err
		}
		file := p.reader.File[i]
		var decompressor Decompressor
		switch file.Comment[0:1] {
		case UseBrotli:
			decompressor = brotliDecompressor
		case UseLZ4:
			decompressor = lz4Decompressor
		case NotToCompress:
			decompressor = nullDecompressor
		}
		return &packedFileEntry{
			file:         file,
			decompressor: decompressor,
			decryptor:    decryptor,
			logicalPath:  path.Clean("/" + path.Join(p.baseBundle.mountPoint, file.Name)),
		}, nil
	} else {
		return nil, nil
	}
}

func (packedBundle) readdir(path string) []FileEntry {
	panic("implement me")
}

func (p *packedBundle) close() {
	if p.closer != nil {
		p.closer.Close()
	}
}

type packedFileEntry struct {
	file         *zip.File
	logicalPath  string
	decompressor Decompressor
	decryptor    Decryptor
}

func (z packedFileEntry) Reader() (io.ReadCloser, error) {
	reader, err := z.file.Open()
	if err != nil {
		return nil, err
	}
	decryptoReader, err := z.decryptor.Decrypto(reader)
	if err != nil {
		return nil, err
	}
	return NewReadCloser(z.decompressor(decryptoReader), reader), nil
}

func (z packedFileEntry) BrotliReader() (io.ReadCloser, error) {
	if z.file.Comment[0:1] == UseBrotli {
		reader, err := z.file.Open()
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

func (z packedFileEntry) Stat() os.FileInfo {
	originalSize, _ := strconv.ParseInt(strings.Split(z.file.Comment[1:], "-")[0], 16, 64)
	return &zipFileInfo{
		name:         z.Name(),
		modTime:      z.file.Modified,
		originalSize: originalSize,
	}
}

func (z packedFileEntry) Name() string {
	return path.Base(z.file.Name)
}

func (z packedFileEntry) Path() string {
	return z.logicalPath
}

func (z packedFileEntry) Etag() string {
	return z.file.Comment[1:]
}

type zipFileInfo struct {
	name         string
	modTime      time.Time
	originalSize int64
}

func (z zipFileInfo) Name() string {
	return z.name
}

func (z zipFileInfo) Size() int64 {
	return z.originalSize
}

func (z zipFileInfo) Mode() os.FileMode {
	return 0444
}
func (z zipFileInfo) ModTime() time.Time {
	return z.modTime
}

func (z zipFileInfo) IsDir() bool {
	return false
}

func (z zipFileInfo) Sys() interface{} {
	return nil
}
