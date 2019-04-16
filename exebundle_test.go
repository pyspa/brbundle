package brbundle_test

import (
	"github.com/shibukawa/brbundle"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestExecutionBundle_Windows_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.exe"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestExecutionBundle_Windows_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.exe"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}

func TestExecutionBundle_Linux_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.linux"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestExecutionBundle_Linux_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.linux"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}

func TestExecutionBundle_Darwin_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.darwin"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestExecutionBundle_Darwin_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionBundle(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.darwin"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}
