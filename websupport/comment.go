package websupport

import (
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func MakeCommentString(compressorFlag, filePath string, info os.FileInfo) string {
	etag := fmt.Sprintf("%x-%x", int(info.Size()), info.ModTime().Unix())
	mimeType := "application/octet-stream"
	if m, err := mimetype.DetectFile(filePath); err == nil {
		mimeType = m.String()
	}
	if strings.HasPrefix(mimeType, "text/plain") {
		if strings.HasSuffix(filePath, ".js") {
			// github.com/gabriel-vasile/mimetype detect js via shebang
			// It fails to detect js files from webpack
			mimeType = strings.Replace(mimeType, "text/plain", "application/javascript", 1)
		} else if strings.HasSuffix(filePath, ".css") {
			// github.com/gabriel-vasile/mimetype doesn't support css
			mimeType = strings.Replace(mimeType, "text/plain", "application/css", 1)
		}
	}
	return fmt.Sprintf("%s,%s,%s", compressorFlag, etag, mimeType)
}
