package brhttp

import (
	"io"
	"net/http"

	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/websupport"
)

type FileSystem struct {
	option brbundle.WebOption
}

func (f FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	file, found, redirectDir := websupport.FindFile(p, f.option)
	if redirectDir {
		http.Redirect(w, r, "./", http.StatusFound)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reader, etag, headers, err := websupport.GetContent(file, f.option, r.Header.Get("Accept-Encoding"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	for _, header := range headers {
		w.Header().Set(header[0], header[1])
	}
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	} else {
		io.Copy(w, reader)
	}
}

func Mount(option ...brbundle.WebOption) *FileSystem {
	return &FileSystem{
		option: websupport.InitOption(option),
	}
}

func MountFunc(option ...brbundle.WebOption) http.HandlerFunc {
	fs := Mount(option...)
	return fs.ServeHTTP
}