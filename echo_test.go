package brbundle_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, nil, handler(c))

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "hello world from root\n", rec.Body.String())
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

	assert.Equal(t, nil, handler(c))
	assert.Equal(t, "br", rec.Header().Get("Content-Encoding"))
	assert.Equal(t, 200, rec.Code)
	reader, err := brotli.NewReader(rec.Body, nil)
	assert.Equal(t, nil, err)
	body, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)

	assert.Equal(t, "hello world from root\n", string(body))
}
