package brbundle_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.pyspa.org/brbundle"
	"go.pyspa.org/brbundle/brhttp"
)

func initServer(rootPath string) *httptest.Server {
	r := brbundle.NewRepository()
	r.RegisterFolder(rootPath)
	return httptest.NewServer(brhttp.Mount(brbundle.WebOption{
		Repository: r,
	}))
}

func TestManifestTestFunction(t *testing.T) {
	ts := initServer("./testdata/result/old-manifest")
	defer ts.Close()
	r, err := http.Get(ts.URL + "/manifest.json")
	if err != nil {
		assert.NotNil(t, err)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		assert.NotNil(t, err)
		return
	}
	t.Log(string(data))
	assert.True(t, strings.Contains(string(data), `"2.3/main/zabbix":`))
}
