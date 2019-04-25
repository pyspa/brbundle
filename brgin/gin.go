package brgin

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/websupport"
)

func Mount(option ...brbundle.WebOption) gin.HandlerFunc {
	o := websupport.InitOption(option)

	return func (c *gin.Context) {
		var p string
		if o.SPAFallback == "" {
			p = c.Param("filepath")
		} else {
			p = c.Request.URL.Path
		}

		file, found, redirectDir := websupport.FindFile(p, o)
		if redirectDir {
			c.Redirect(http.StatusFound, "./")
			return
		} else if !found {
			c.Status(http.StatusNotFound)
			return
		}

		reader, etag, headers, err := websupport.GetContent(file, o, string(c.GetHeader("Accept-Encoding")))
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}
		defer reader.Close()

		for _, header := range headers {
			c.Header(header[0], header[1])
		}
		if string(c.GetHeader("If-None-Match")) == etag {
			c.Status(http.StatusNotModified)
			return
		} else {
			defer reader.Close()
			io.Copy(c.Writer, reader)
		}
	}
}
