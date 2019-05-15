package websupport

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"go.pyspa.org/brbundle"
)

const DefaultQValue = 1.0

// based on https://github.com/NYTimes/gziphandler/blob/master/gzip.go
func HasSupportBrotli(s string) bool {
	var e []string

	for _, ss := range strings.Split(s, ",") {
		coding, qvalue, err := parseCoding(ss)
		if err != nil {
			e = append(e, err.Error())
		} else {
			if coding == "br" {
				return qvalue > 0.0
			}
		}
	}
	return false
}

// parseCoding parses a single conding (content-coding with an optional qvalue),
// as might appear in an Accept-Encoding header. It attempts to forgive minor
// formatting errors.
func parseCoding(s string) (coding string, qvalue float64, err error) {
	for n, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		qvalue = DefaultQValue

		if n == 0 {
			coding = strings.ToLower(part)
		} else if strings.HasPrefix(part, "q=") {
			qvalue, err = strconv.ParseFloat(strings.TrimPrefix(part, "q="), 64)

			if qvalue < 0.0 {
				qvalue = 0.0
			} else if qvalue > 1.0 {
				qvalue = 1.0
			}
		}
	}

	if coding == "" {
		err = fmt.Errorf("empty content-coding")
	}

	return
}

func FindFile(p string, o brbundle.WebOption) (file brbundle.FileEntry, found, redirectDir bool) {
	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	checkDir := false
	if strings.HasSuffix(p, "/") && o.DirectoryIndex != "" {
		p = p + o.DirectoryIndex
		checkDir = true
	}

	var err error

	file, err = o.Repository.Find(p)
	if file == nil {
		if o.DirectoryIndex != "" && !checkDir {
			if p == "" {
				p = o.DirectoryIndex
			} else {
				p = p + "/" + o.DirectoryIndex
			}
			file, err = o.Repository.Find(p)
			if file != nil {
				found = true
				return
			}
		}
		if o.SPAFallback != "" {
			file, err = o.Repository.Find(o.SPAFallback)
			if err != nil {
				found = true
				return
			}
		} else {
			found = false
			return
		}
	}
	found = true
	return
}

func GetContent(file brbundle.FileEntry, o brbundle.WebOption, acceptEncoding string) (reader io.ReadCloser, etag string, headers [][2]string, err error) {
	headers = make([][2]string, 3, 5)
	etag, contentType := file.EtagAndContentType()
	headers[0] = [2]string{"ETag", etag}
	headers[1] = [2]string{"Content-Type", contentType}

	brreader, _ := file.BrotliReader()
	if brreader != nil && HasSupportBrotli(acceptEncoding) {
		headers[2] = [2]string{"Content-Length", strconv.FormatInt(file.CompressedSize(), 10)}
		headers = append(headers, [2]string{"Content-Encoding", "br"})
		reader = brreader
	} else {
		headers[2] = [2]string{"Content-Length", strconv.FormatInt(file.Stat().Size(), 10)}
		reader, err = file.Reader()
	}
	if o.MaxAge != time.Duration(0) {
		headers = append(headers, [2]string{"Cache-Control", fmt.Sprintf("max-age=%d", int(o.MaxAge.Seconds()))})
	}

	return

}

func InitOption(options []brbundle.WebOption) brbundle.WebOption {
	var o brbundle.WebOption
	if len(options) > 0 {
		o = options[0]
	} else {
		o.DirectoryIndex = "index.html"
	}
	if o.Repository == nil {
		o.Repository = brbundle.DefaultRepository
	}
	return o
}
