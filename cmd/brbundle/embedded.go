package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

const embeddedSourceTemplate = `[[.BuildTag]]package [[.PackageName]]

import (
    "go.pyspa.org/brbundle"
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
	buildTag     string
}

func (c Context) Content() string {
	return formatContent(c.zipContent.Bytes(), 70)
}

func (c Context) BundleName() string {
	return fmt.Sprintf("%#v", c.bundleName)
}

func (c Context) BuildTag() string {
	if c.buildTag == "" {
		return ""
	}
	return fmt.Sprintf("// +build %s\n\n", c.buildTag)
}

func embedded(brotli bool, encryptionKey []byte, buildTag, packageName string, destFile *os.File, srcDirPath, dirPrefix, bundleName string, date *time.Time) error {
	var zipContent bytes.Buffer
	packedBundle(brotli, encryptionKey, buildTag, &zipContent, srcDirPath, dirPrefix, "Embedded File", date)

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
		buildTag:     buildTag,
		zipContent:   &zipContent,
	}

	color.HiGreen("Writing %s\n", destFile.Name())
	t := template.Must(template.New("embedded").Delims("[[", "]]").Parse(embeddedSourceTemplate))
	var source bytes.Buffer
	err = t.Execute(&source, *context)
	if err != nil {
		panic(err)
	}
	destFile.Write(source.Bytes())
	color.Cyan("\nDone\n\n")
	return nil
}

func splitByte(src []byte, length int) []string {
	var result []string

	builder := &strings.Builder{}
	for _, c := range src {
		printable := 32 <= c && c <= 126 && c != 34 && c != 92
		if printable {
			builder.WriteByte(c)
		} else if c == 10 {
			builder.WriteString(`\n`)
		} else if c == 34 {
			builder.WriteString(`\"`)
		} else if c == 92 {
			builder.WriteString(`\\`)
		} else {
			fmt.Fprintf(builder, "\\x%02x", c)
		}
		if builder.Len() > length {
			result = append(result, builder.String())
			builder.Reset()
		}
	}
	if builder.Len() > 0 {
		result = append(result, builder.String())
	}
	if len(result) == 0 {
		result = append(result, "")
	}

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
