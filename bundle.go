package brbundle

import (
	"os"
	"io"
	"fmt"
	"sort"
)

type FileEntry interface{
	Reader() (io.Reader, error)
	BrotliReader() (io.Reader, error)
	Stat() os.FileInfo
	Name() string
	Path() string
}

type FilePod interface {
	Find(path string) FileEntry
	Readdir(path string) []FileEntry
	Close() error
}

type Bundle struct {
	pods []FilePod
}

func NewBundle(pods ...FilePod) *Bundle {
	return &Bundle{pods}
}

func (b Bundle) Find(path string) (FileEntry, error) {
	for _, pod := range b.pods {
		entry := pod.Find(path)
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
	for _, pod := range b.pods {
		entries := pod.Readdir(path)
		if entries != nil {
			found = true
			for _, entry := range entries {
				if foundFiles[entry.Name()] == nil {
					foundFiles[entry.Name()] = entry;
					fileNames = append(fileNames, entry.Name())
				}
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("Can't read the dir: %", path)
	}
	sort.Strings(fileNames)
	result := make([]FileEntry, len(fileNames))
	for i, fileName := range fileNames {
		result[i] = foundFiles[fileName]
	}
	return result, nil
}

func (b *Bundle) Close(deletePod FilePod) error {
	var pods []FilePod
	if len(b.pods) > 1 {
		pods = make([]FilePod, 0, len(b.pods) - 1)
	}

	var err error

	for _, pod := range b.pods {
		if pod != deletePod {
			pods = append(pods, pod)
		} else {
			err = pod.Close()
		}
	}
	return err
}