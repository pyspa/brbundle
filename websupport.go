package brbundle

import (
	"strings"
	"time"
)

type WebOption struct {
	Repository     *Repository
	SPAFallback    string
	MaxAge         time.Duration
	DirectoryIndex string
}

func ParseCommentString(comment string) (compressorFlag, etag, contentType string) {
	result := strings.SplitN(comment, ",", 3)
	compressorFlag = result[0]
	etag = result[1]
	contentType = result[2]
	return
}
