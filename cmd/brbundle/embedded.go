package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
	"go/format"

	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/fatih/color"
	"github.com/shibukawa/brbundle"
	"path"
	"sync"
	"sort"
)

const embeddedSourceTemplate = `
package [[.PackageName]]

import (
    "time"
    "github.com/shibukawa/brbundle"
)

[[range .Files]][[if ne .UniqueName "nil"]]
var [[.UniqueName]] = [[.Content]]
[[end]][[end]]

// [[.ConstantName]] returns content pod for brbundle FileSystem
var [[.ConstantName]] = brbundle.MustEmbeddedPod([[.Decompressor]], [[.Deencryptor]], [[.Dirs]], map[string]*brbundle.Entry{[[range .Files]]
"[[.Path]]": &brbundle.Entry{
    Path: "[[.Path]]",
    FileMode: [[.FileMode]],
	OriginalSize: [[.OriginalSize]],
	Mtime: [[.Mtime]],
    Data: [[.UniqueName]],
    ETag: "[[.ETag]]",
},
[[end]]})

`

type File struct {
	path       string
	content    []byte
	ETag       string
	info       os.FileInfo
	UniqueName string
}

func (f File) Path() string {
	return filepath.ToSlash(f.path)
}

func (f File) FileMode() string {
	return fmt.Sprintf("%#o", f.info.Mode())
}

func (f File) OriginalSize() int64 {
	return f.info.Size()
}

func (f File) Mtime() string {
	return fmt.Sprintf("time.Unix(%#v, %#v)", f.info.ModTime().Unix(), f.info.ModTime().UnixNano())
}

func (f File) Content() string {
	return fmt.Sprintf("[]byte(%#v)", string(f.content))
}

type Context struct {
	compression  brbundle.CompressionType
	encryption   brbundle.EncryptionType
	PackageName  string
	VariableName string
	files        []*File
	dirs         map[string][]string
	lock         sync.Mutex
}

func (c Context) Decompressor() string {
	return c.compression.FunctionName()
}

func (c Context) Deencryptor() string {
	return c.encryption.FunctionName()
}

func (c Context) Dirs() string {
	keys := make([]string, 0, len(c.dirs))
	for key := range c.dirs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "map[string][]string{\n")
	for _, key := range keys {
		fmt.Fprintf(&buffer, "\t%#v: []string{", key)
		var files = c.dirs[key]
		sort.Strings(files)
		for j, file := range files {
			if j != 0 {
				fmt.Fprintf(&buffer, ", ")
			}
			fmt.Fprintf(&buffer, `"%s"`, file)
		}
		fmt.Fprintf(&buffer, "},\n")
	}
	fmt.Fprintf(&buffer, "}")
	return buffer.String()
}

type Files []*File

// 以下インタフェースを満たす

func (f Files) Len() int {
	return len(f)
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Files) Less(i, j int) bool {
	return f[i].path < f[j].path
}
func (c Context) Files() []*File {
	var files Files = c.files
	sort.Sort(files)
	return files
}

func (c *Context) AddFile(entry Entry, etag string, writerTo io.WriterTo) {
	var buffer bytes.Buffer
	hash := sha1.New()
	writer := io.MultiWriter(&buffer, hash)
	writerTo.WriteTo(writer)

	c.lock.Lock()
	defer c.lock.Unlock()

	filePath := "/" + filepath.ToSlash(entry.Path)
	c.files = append(c.files, &File{
		path:       filePath,
		content:    buffer.Bytes(),
		ETag:       etag,
		info:       entry.Info,
		UniqueName: fmt.Sprintf("_%s%x", c.VariableName, hash.Sum(nil)),
	})
	dir := path.Dir(filePath)
	c.dirs[dir] = append(c.dirs[dir], path.Base(filePath))
}

func (c *Context) AddDir(entry Entry) {
	filePath := "/" + filepath.ToSlash(entry.Path)
	c.files = append(c.files, &File{
		path:       filePath,
		content:    nil,
		ETag:       "",
		info:       entry.Info,
		UniqueName: "nil",
	})
}

func embeddedWorker(compressor *Compressor, encryptor *Encryptor, context *Context, srcDirPath string, jobs <-chan Entry, wait chan<- struct{}) {
	for entry := range jobs {
		err := processInput(compressor, encryptor, srcDirPath, entry, func(writerTo io.WriterTo, etag string) {
			context.AddFile(entry, etag, writerTo)
		})

		if err != nil {
			continue
		}
	}
	wait <- struct{}{}
}

func embedded(ctype brbundle.CompressionType, etype brbundle.EncryptionType, ekey []byte, packageName, variableName string, destFile *os.File, srcDirPath string) {
	color.Cyan("Embedded File Mode (Compression: %s, Encyrption: %s)\n\n", ctype, etype)

	paths, dirs, ignored := Traverse(srcDirPath)

	defer destFile.Close()

	context := &Context{
		compression:  ctype,
		encryption:   etype,
		PackageName:  packageName,
		VariableName: variableName,
		files:        make([]*File, 0, len(paths)),
		dirs:         make(map[string][]string),
	}

	for _, dir := range dirs {
		context.AddDir(dir)
	}

	wait := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go embeddedWorker(NewCompressor(ctype), NewEncryptor(etype, ekey), context, srcDirPath, paths, wait)
	}

	close(paths)

	for i := 0; i < runtime.NumCPU(); i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}

	color.HiGreen("Writing %s\n", destFile.Name())
	t := template.Must(template.New("embedded").Delims("[[", "]]").Parse(embeddedSourceTemplate))
	var source bytes.Buffer
	t.Execute(&source, *context)
	formattedSource, err := format.Source(source.Bytes())
	if err != nil {
		fmt.Println(err.Error())
		destFile.Write(source.Bytes())
	} else {
		destFile.Write(formattedSource)
	}
	color.Cyan("\nDone\n\n")
}
