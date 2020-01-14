package websupport

import (
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func MakeCommentString(compressorFlag, filePath string, info os.FileInfo) string {
	etag := fmt.Sprintf("%x-%x", int(info.Size()), info.ModTime().Unix())
	mimeType := "application/octet-stream"
	if m, err := mimetype.DetectFile(filePath); err != nil {
		mimeType = m.String()
	}
	return fmt.Sprintf("%s,%s,%s", compressorFlag, etag, mimeType)
}
