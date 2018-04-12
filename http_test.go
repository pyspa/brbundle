package brbundle_test

import (
	"github.com/ToQoz/gopwt/assert"
	"github.com/shibukawa/brbundle"
	"io/ioutil"
	"testing"
	"net/http/httptest"
	"net/http"
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