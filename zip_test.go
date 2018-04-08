package brbundle_test

import (
	"github.com/ToQoz/gopwt/assert"
	"github.com/shibukawa/brbundle"
	"io/ioutil"
	"testing"
)

func TestZipRawNoCrypto_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.NullDecompressor(), brbundle.NullDecryptor(), "./testdata/raw-nocrypto.zip"))
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

func TestZipRawNoCrypto_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.NullDecompressor(), brbundle.NullDecryptor(), "./testdata/raw-nocrypto.zip"))
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

func TestZipBrotliAES_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.AESDecryptor(encryptKey), "./testdata/br-aes.zip"))
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

func TestZipBrotliAES_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.AESDecryptor(encryptKey),"./testdata/br-aes.zip"))
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

func TestZipLZ4ChaCha_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.LZ4Decompressor(), brbundle.ChaChaDecryptor(encryptKey), "./testdata/lz4-chacha.zip"))
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

func TestZipLZ4ChaCha_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.LZ4Decompressor(), brbundle.ChaChaDecryptor(encryptKey), "./testdata/lz4-chacha.zip"))
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
