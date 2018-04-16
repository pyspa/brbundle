package brbundle_test

import (
	"io/ioutil"
	"testing"
	"net/http/httptest"
	"net/http"

	"github.com/ToQoz/gopwt/assert"
	"github.com/shibukawa/brbundle"
	"github.com/dsnet/compress/brotli"
)


func TestNewFileSystem_NoBrotli(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))

	s := httptest.NewServer(brbundle.MountBundle("/static", bundle))
	defer s.Close()

	res, err := http.Get(s.URL + "/static/rootfile.txt")
	assert.OK(t, err == nil)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.OK(t, err == nil)

	assert.OK(t, string(body) == "hello world from root\n")
}

func TestNewFileSystem_Brotli(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))

	s := httptest.NewServer(brbundle.MountBundle("/static", bundle))
	defer s.Close()

	request, err := http.NewRequest("GET", s.URL + "/static/rootfile.txt", nil)
	assert.OK(t, err == nil)
	request.Header.Add("Accept-Encoding", "br")
	res, err := http.DefaultClient.Do(request)
	assert.OK(t, err == nil)

	assert.OK(t, res.Header.Get("Content-Encoding") == "br")

	defer res.Body.Close()
	reader, err := brotli.NewReader(res.Body, nil)
	assert.OK(t, err == nil)
	body, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)

	assert.OK(t, string(body) == "hello world from root\n")
}