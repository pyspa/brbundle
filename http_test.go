package brbundle_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/dsnet/compress/brotli"
	"github.com/shibukawa/brbundle"
)

func TestNewFileSystem_NoBrotli(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))

	s := httptest.NewServer(brbundle.ServerMount("/static", bundle))
	defer s.Close()

	res, err := http.Get(s.URL + "/static/rootfile.txt")
	assert.Equal(t, nil, err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)

	assert.Equal(t, "hello world from root\n", string(body))
}

func TestNewFileSystem_Brotli(t *testing.T) {
	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))

	s := httptest.NewServer(brbundle.ServerMount("/static", bundle))
	defer s.Close()

	request, err := http.NewRequest("GET", s.URL+"/static/rootfile.txt", nil)
	assert.Equal(t, nil, err)
	request.Header.Add("Accept-Encoding", "br")
	res, err := http.DefaultClient.Do(request)
	assert.Equal(t, nil, err)

	assert.Equal(t, "br", res.Header.Get("Content-Encoding"))

	defer res.Body.Close()
	reader, err := brotli.NewReader(res.Body, nil)
	assert.Equal(t, nil, err)
	body, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)

	assert.Equal(t, "hello world from root\n", string(body))
}
