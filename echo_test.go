package brbundle_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/dsnet/compress/brotli"
	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle"
)

func TestNewFileSystem_ForEcho_NoBrotli(t *testing.T) {
	e := echo.New()

	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))
	handler := brbundle.EchoMount("/static", bundle)

	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/static/rootfile.txt")

	assert.OK(t, handler(c) == nil)

	assert.OK(t, rec.Code == 200)
	assert.OK(t, rec.Body.String() == "hello world from root\n")
}

func TestNewFileSystem_ForEcho_Brotli(t *testing.T) {
	e := echo.New()

	var bundle = brbundle.NewBundle(brbundle.MustZipPod(brbundle.BrotliDecompressor(), brbundle.NullDecryptor(), "./testdata/br-nocrypto.zip"))
	handler := brbundle.EchoMount("/static", bundle)

	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Add("Accept-Encoding", "br")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/static/rootfile.txt")
	//c.Set("Accept-Encoding", "br")

	assert.OK(t, handler(c) == nil)
	assert.OK(t, rec.Header().Get("Content-Encoding") == "br")
	assert.OK(t, rec.Code == 200)
	reader, err := brotli.NewReader(rec.Body, nil)
	assert.OK(t, err == nil)
	body, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)

	assert.OK(t, string(body) == "hello world from root\n")
}
