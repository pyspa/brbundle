package brbundle

import (
	"fmt"
	"io"
	"os"
	"sort"
)

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
}

type FileBundle interface {
	Find(path string) FileEntry
	Readdir(path string) []FileEntry
	Close() error
}

type Bundle struct {
	bundles []FileBundle
}

func NewBundle(bundles ...FileBundle) *Bundle {
	return &Bundle{bundles}
}

func (b *Bundle) AddBundle(bundle FileBundle) {
	b.bundles = append(b.bundles, bundle)
}

func (b Bundle) Find(path string) (FileEntry, error) {
	for _, bundle := range b.bundles {
		entry := bundle.Find(path)
		if entry != nil {
			return entry, nil
		}
	}
	return nil, fmt.Errorf("Can't read the file: %s", path)
}

func (b Bundle) Readdir(path string) ([]FileEntry, error) {
	var foundFiles = make(map[string]FileEntry)
	var fileNames []string
	var found = false
	for _, bundle := range b.bundles {
		entries := bundle.Readdir(path)
		if entries != nil {
			found = true
			for _, entry := range entries {
				if foundFiles[entry.Name()] == nil {
					foundFiles[entry.Name()] = entry
					fileNames = append(fileNames, entry.Name())
				}
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("Can't read the dir: %s", path)
	}
	sort.Strings(fileNames)
	result := make([]FileEntry, len(fileNames))
	for i, fileName := range fileNames {
		result[i] = foundFiles[fileName]
	}
	return result, nil
}

func (b *Bundle) Close(deleteBundle FileBundle) error {
	var bundles []FileBundle
	if len(b.bundles) > 1 {
		bundles = make([]FileBundle, 0, len(b.bundles)-1)
	}

	var err error

	for _, bundle := range b.bundles {
		if bundle != deleteBundle {
			bundles = append(bundles, bundle)
		} else {
			err = bundle.Close()
		}
	}
	return err
}
