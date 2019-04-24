package websupport

import (
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func MakeCommentString(compressorFlag, filePath string, info os.FileInfo) string {
	etag := fmt.Sprintf("%x-%x", int(info.Size()), info.ModTime().Unix())
	mimeType, _, _ := mimetype.DetectFile(filePath)
	return fmt.Sprintf("%s,%s,%s", compressorFlag, etag, mimeType)
}
