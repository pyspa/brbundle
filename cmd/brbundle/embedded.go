package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go/format"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

const embeddedSourceTemplate = `
package [[.PackageName]]

import (
    "github.com/shibukawa/brbundle"
)

var [[.VariableName]] = [[.Content]]

func init() {
    brbundle.RegisterEmbeddedBundle([[.VariableName]], [[.BundleName]])
}
`

type Context struct {
	PackageName  string
	zipContent   *bytes.Buffer
	VariableName string
	bundleName   string
	lock         sync.Mutex
}

func (c Context) Content() string {
	return formatContent(c.zipContent.Bytes(), 70)
}

func (c Context) BundleName() string {
	return fmt.Sprintf("%#v", c.bundleName)
}

func embedded(brotli bool, encryptionKey []byte, packageName string, destFile *os.File, srcDirPath, dirPrefix, bundleName string, date *time.Time) error {
	var zipContent bytes.Buffer
	packedBundle(brotli, encryptionKey, &zipContent, srcDirPath, dirPrefix, "Embedded File", date)

	_, err := NewEncryptor(encryptionKey)
	if err != nil {
		return errors.New("Can't create encryptor")
	}

	defer destFile.Close()

	h := md5.New()
	h.Write(zipContent.Bytes())

	context := &Context{
		PackageName:  packageName,
		VariableName: fmt.Sprintf("bundle_%s", hex.EncodeToString(h.Sum(nil))),
		bundleName:   bundleName,
		zipContent:   &zipContent,
	}

	color.HiGreen("Writing %s\n", destFile.Name())
	t := template.Must(template.New("embedded").Delims("[[", "]]").Parse(embeddedSourceTemplate))
	var source bytes.Buffer
	err = t.Execute(&source, *context)
	if err != nil {
		panic(err)
	}
	formattedSource, err := format.Source(source.Bytes())
	if err != nil {
		destFile.Write(source.Bytes())
	} else {
		destFile.Write(formattedSource)
	}
	color.Cyan("\nDone\n\n")
	return nil
}

func splitByte(src []byte, length int) []string {
	var result []string

	str := fmt.Sprintf("%#v", string(src))
	str = str[1 : len(str)-1]
	start := 0
	for i := 0; i < len(str)-3; i++ {
		if i-start > length {
			result = append(result, str[start:i])
			start = i
		}
		if str[i:i+2] == `\x` {
			i += 3
		} else if str[i:i+2] == `\u` {
			i += 5
		} else if str[i:i+1] == `\` {
			i += 1
		}
	}
	result = append(result, str[start:])

	return result
}

func formatContent(src []byte, length int) string {
	lines := splitByte(src, length)
	quoted := make([]string, len(lines))
	for i, line := range lines {
		quoted[i] = fmt.Sprintf("\"%s\"", line)
	}
	switch len(quoted) {
	case 0:
		return "[]byte(\"\")"
	case 1:
		return fmt.Sprintf("[]byte(%s)", quoted[0])
	default:
		return fmt.Sprintf("[]byte(\n\t%s)", strings.Join(quoted, " +\n\t"))
	}
}
