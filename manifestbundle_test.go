package brbundle_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestManifestInitialDownload(t *testing.T) {
	ts := initServer("./testdata/result/old-manifest")
	defer ts.Close()

	r := brbundle.NewRepository()
	p, err := r.RegisterRemoteManifest(ts.URL, brbundle.Option{
		TempFolder:          filepath.Join(os.TempDir(), "read-test-1"),
		ResetDownloadFolder: true,
	})
	if err != nil {
		assert.Nil(t, err)
		return
	}
	// At first time, brbundle downloads all entries.
	assert.Equal(t, 17, len(p.DownloadFiles()))
	assert.Equal(t, 0, len(p.DeleteFiles()))
	assert.Equal(t, 0, len(p.KeepFiles()))

	p.Wait()
	e, err := r.Find("2.3/main/zabbix/CVE-2013-5743.json")
	if err != nil {
		assert.Nil(t, err)
		return
	}
	body, _ := e.ReadAll()
	assert.True(t, strings.Contains(string(body), `"IssueID": 227`))
}

func TestManifestUpdate(t *testing.T) {
	{
		ts := initServer("./testdata/result/old-manifest")
		defer ts.Close()

		r := brbundle.NewRepository()
		p, err := r.RegisterRemoteManifest(ts.URL, brbundle.Option{
			TempFolder:          filepath.Join(os.TempDir(), "read-test-2"),
			ResetDownloadFolder: true,
		})
		if err != nil {
			assert.Nil(t, err)
			return
		}

		p.Wait()

	}

	{
		ts := initServer("./testdata/result/new-manifest")
		defer ts.Close()

		r := brbundle.NewRepository()
		p, err := r.RegisterRemoteManifest(ts.URL, brbundle.Option{
			TempFolder:          filepath.Join(os.TempDir(), "read-test-2"),
			ResetDownloadFolder: false,
		})
		if err != nil {
			assert.Nil(t, err)
			return
		}

		// At second time, brbundle downloads new&updated entries only.
		assert.Equal(t, 4, len(p.DownloadFiles()))
		assert.Equal(t, 1, len(p.DeleteFiles()))
		assert.Equal(t, 16, len(p.KeepFiles()))

		p.Wait()

		e, err := r.Find("2.3/main/zabbix/CVE-2013-5743.json")
		if err != nil {
			assert.Nil(t, err)
			return
		}
		body, _ := e.ReadAll()
		assert.True(t, strings.Contains(string(body), `"IssueID": 227`))
	}
}
