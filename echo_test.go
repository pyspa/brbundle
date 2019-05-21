package brbundle_test

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dsnet/compress/brotli"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.pyspa.org/brbundle"
	"go.pyspa.org/brbundle/brecho"
)

func TestEchoMount_NoBrotli(t *testing.T) {
	repo := initRepo()

	e := echo.New()

	handler := brecho.Mount(brbundle.WebOption{
		Repository: repo,
	})

	req := httptest.NewRequest(echo.GET, "/rootfile.txt", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/rootfile.txt")

	assert.Equal(t, nil, handler(c))
	assert.Equal(t, 200, rec.Code)
	assert.True(t, strings.HasPrefix(string(rec.Body.Bytes()), "hello world from root\n"))
}

func TestEchoMount_Brotli(t *testing.T) {
	repo := initRepo()

	e := echo.New()

	handler := brecho.Mount(brbundle.WebOption{
		Repository: repo,
	})

	req := httptest.NewRequest(echo.GET, "/rootfile.txt", nil)
	req.Header.Add("Accept-Encoding", "br")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/rootfile.txt")

	assert.Equal(t, nil, handler(c))
	assert.Equal(t, "br", rec.Header().Get("Content-Encoding"))
	assert.Equal(t, 200, rec.Code)
	reader, err := brotli.NewReader(rec.Body, nil)
	assert.Equal(t, nil, err)
	body, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)

	assert.True(t, strings.HasPrefix(string(body), "hello world from root\n"))
}

func TestEchoMount_SPAOption(t *testing.T) {
	repo := initRepo()

	e := echo.New()

	handler := brecho.Mount(brbundle.WebOption{
		Repository:  repo,
		SPAFallback: "index.html",
	})

	req := httptest.NewRequest(echo.GET, "/profile", nil)
	req.Header.Add("Accept-Encoding", "br")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	// undefined
	c.SetPath("/profile")

	assert.Equal(t, nil, handler(c))
	assert.Equal(t, 200, rec.Code)
	assert.True(t, strings.Contains(string(rec.Body.Bytes()), "<body>"))
}
