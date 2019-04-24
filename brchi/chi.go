package brchi

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/websupport"
)

func Mount(option ...brbundle.WebOption) http.HandlerFunc {
	o := websupport.InitOption(option)

	return func (w http.ResponseWriter, r *http.Request) {
		p := chi.URLParam(r, "*")

		file, found, redirectDir := websupport.FindFile(p, o)
		if redirectDir {
			http.Redirect(w, r, "./", http.StatusFound)
			return
		} else if !found {
			w.WriteHeader(404)
			return
		}

		reader, etag, headers, err := websupport.GetContent(file, o, r.Header.Get("Accept-Encoding"))
		if err != nil {
			w.WriteHeader(500)
		}
		defer reader.Close()

		for _, header := range headers {
			w.Header().Set(header[0], header[1])
		}
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			return
		} else {
			defer reader.Close()
			io.Copy(w, reader)
		}
	}
}
