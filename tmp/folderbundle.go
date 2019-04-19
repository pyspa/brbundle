// +build !js

package brbundle

type FSBundle struct {
	path string
}

func (FSBundle) Find(path string) (FileEntry, bool) {
	panic("implement me")
}

func (FSBundle) Readdir(path string) []FileEntry {
	panic("implement me")
}

func (FSBundle) Close() error {
	panic("implement me")
}

func (FSBundle) SetDecryptionKey(key []byte) {
	panic("implement me")
}