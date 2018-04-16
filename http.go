package brbundle

import (
	"io"
	"net/http"
	"strings"
)

type FileSystem struct {
	path   string
	bundle *Bundle
}

func (f FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newPath := strings.TrimPrefix(r.URL.Path, f.path)
	if newPath == r.URL.Path {
		w.WriteHeader(404)
		return
	}

	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	file, err := f.bundle.Find(newPath)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if supportBrotli(r.Header.Get("Accept-Encoding")) {
		reader, err := file.BrotliReader()
		w.Header().Set("Content-Encoding", "br")
		defer reader.Close()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		io.Copy(w, reader)
	} else {
		reader, err := file.Reader()
		defer reader.Close()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		io.Copy(w, reader)
	}
}

func ServerMount(path string, bundle *Bundle) *FileSystem {
	return &FileSystem{
		path:   path,
		bundle: bundle,
	}
}
