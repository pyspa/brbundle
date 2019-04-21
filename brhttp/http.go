package brhttp

import (
	"io"
	"net/http"
	"strings"
	"github.com/shibukawa/brbundle"
)

type FileSystem struct {
	path   string
	option brbundle.WebOption
}

func (f FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newPath := strings.TrimPrefix(r.URL.Path, f.path)
	if newPath == r.URL.Path {
		w.WriteHeader(404)
		return
	}

	if strings.HasPrefix(newPath, "/") {
		newPath = newPath[1:]
	}

	file, err := f.option.Repository.Find(newPath)
	if err != nil {
		if f.option.SPAFallback != "" {
			file, err = f.option.Repository.Find(f.option.SPAFallback)
			if err != nil {
				w.WriteHeader(404)
				return
			}
		} else {
			w.WriteHeader(404)
			return
		}
	}
	brreader, err := file.BrotliReader()
	if brreader != nil && brbundle.HasSupportBrotli(r.Header.Get("Accept-Encoding")) {
		w.Header().Set("Content-Encoding", "br")
		defer brreader.Close()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		io.Copy(w, brreader)
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

func Mount(path string, option ...brbundle.WebOption) *FileSystem {
	fs := &FileSystem{
		path:   path,
	}
	if len(option) > 0 {
		fs.option = option[0]
	}
	if fs.option.Repository == nil {
		fs.option.Repository = brbundle.DefaultRepository
	}
	return fs
}
