package brbundle_test

import (
	"github.com/ToQoz/gopwt/assert"
	"github.com/shibukawa/brbundle"
	"io/ioutil"
	"testing"
)

func TestExecutionPod_Windows_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.exe"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "rootfile.txt")
	assert.OK(t, entry.Path() == "/rootfile.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from root\n")
}

func TestExecutionPod_Windows_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.exe"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "file-in-subfolder.txt")
	assert.OK(t, entry.Path() == "/subfolder/file-in-subfolder.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from subfolder\n")
}

func TestExecutionPod_Linux_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.linux"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "rootfile.txt")
	assert.OK(t, entry.Path() == "/rootfile.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from root\n")
}

func TestExecutionPod_Linux_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.linux"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "file-in-subfolder.txt")
	assert.OK(t, entry.Path() == "/subfolder/file-in-subfolder.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from subfolder\n")
}

func TestExecutionPod_Darwin_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.darwin"))
	entry, err := bundle.Find("/rootfile.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "rootfile.txt")
	assert.OK(t, entry.Path() == "/rootfile.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from root\n")
}

func TestExecutionPod_Darwin_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustExecutionPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/testexe/testexe.darwin"))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "file-in-subfolder.txt")
	assert.OK(t, entry.Path() == "/subfolder/file-in-subfolder.txt")
	reader, err := entry.Reader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from subfolder\n")
}
