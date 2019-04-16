package brbundle_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/shibukawa/brbundle"
	"io/ioutil"
	"testing"
)

func TestZipRawNoCrypto_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.NullDecompressor(), brbundle.NullDecryptor(), "./testdata/raw-nocrypto.zip"))
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

func TestZipRawNoCrypto_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.NullDecompressor(), brbundle.NullDecryptor(), "./testdata/raw-nocrypto.zip"))
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

func TestZipBrotliAES_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.BrotliDecompressor(), brbundle.AESDecryptor(encryptKey), "./testdata/br-aes.zip"))
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

func TestZipBrotliAES_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.BrotliDecompressor(), brbundle.AESDecryptor(encryptKey), "./testdata/br-aes.zip"))
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

func TestZipLZ4ChaCha_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.LZ4Decompressor(), brbundle.ChaChaDecryptor(encryptKey), "./testdata/lz4-chacha.zip"))
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

func TestZipLZ4ChaCha_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(brbundle.MustZipBundle(brbundle.LZ4Decompressor(), brbundle.ChaChaDecryptor(encryptKey), "./testdata/lz4-chacha.zip"))
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
