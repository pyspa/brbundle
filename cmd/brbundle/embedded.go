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
	"sync"
	"text/template"
	"time"
)

const embeddedSourceTemplate = `
package [[.PackageName]]

import (
    "github.com/shibukawa/brbundle"
)

// [[.VariableName]] returns content bundle for brbundle FileSystem
var [[.VariableName]] = [[.Content]]

func init() {
    brbundle.MustReadBytes([[.VariableName]], brbundle.OrderEmbedded)
}
`

type Context struct {
	PackageName  string
	zipContent   *bytes.Buffer
	VariableName string
	lock         sync.Mutex
}

func (c Context) Content() string {
	return fmt.Sprintf("[]byte(%#v)", string(c.zipContent.Bytes()))
}

func embedded(brotli bool, encryptionKey []byte, packageName string, destFile *os.File, srcDirPath string, date *time.Time) error {
	var zipContent bytes.Buffer
	zipBundle(brotli, encryptionKey, &zipContent, srcDirPath, "Embedded File", date)

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
