// +build !js

package brbundle

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (p *packedFileEntry) GetLocalPath() (string, error) {
	etag, _ := p.EtagAndContentType()
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s__%s", etag, p.Name()))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, p.Stat().Mode()|0200)
	if err != nil {
		return "", err
	}
	defer f.Close()
	r, err := p.Reader()
	if err != nil {
		return "", err
	}
	defer r.Close()
	io.Copy(f, r)
	return path, nil
}
