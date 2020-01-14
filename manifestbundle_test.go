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

func TestManifest_Dirs(t *testing.T) {
	r := brbundle.NewRepository()
	{
		ts := initServer("./testdata/result/old-manifest")
		defer ts.Close()

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
	dirs := r.Dirs()
	assert.Equal(t, 18, len(dirs))
}

func TestRepository_Manifest(t *testing.T) {
	testcases := []struct {
		Name       string
		MountPoint string
		Dirs       []string
		TestDir    string
		FilesInDir []string
	}{
		{
			"no mount point",
			"",
			[]string{
				"2.2/main/apache2",
				"2.2/main/curl",
				"2.2/main/linux-vserver",
				"2.2/main/openldap",
				"2.2/main/openssl",
				"2.2/main/openvpn",
				"2.2/main/php",
				"2.2/main/samba",
				"2.3/main/apache2",
				"2.3/main/bind",
				"2.3/main/cgit",
				"2.3/main/curl",
				"2.3/main/linux-vserver",
				"2.3/main/openssl",
				"2.3/main/openvpn",
				"2.3/main/xen",
				"2.3/main/zabbix",
			},
			"2.2/main/apache2",
			[]string{
				"CVE-2011-3368.json",
				"CVE-2011-3607.json",
				"CVE-2012-0021.json",
				"CVE-2012-0031.json",
				"CVE-2012-0053.json",
			},
		},
		{
			"has mount point",
			"mount",
			[]string{
				"mount/2.2/main/apache2",
				"mount/2.2/main/curl",
				"mount/2.2/main/linux-vserver",
				"mount/2.2/main/openldap",
				"mount/2.2/main/openssl",
				"mount/2.2/main/openvpn",
				"mount/2.2/main/php",
				"mount/2.2/main/samba",
				"mount/2.3/main/apache2",
				"mount/2.3/main/bind",
				"mount/2.3/main/cgit",
				"mount/2.3/main/curl",
				"mount/2.3/main/linux-vserver",
				"mount/2.3/main/openssl",
				"mount/2.3/main/openvpn",
				"mount/2.3/main/xen",
				"mount/2.3/main/zabbix",
			},
			"mount/2.2/main/apache2",
			[]string{
				"CVE-2011-3368.json",
				"CVE-2011-3607.json",
				"CVE-2012-0021.json",
				"CVE-2012-0031.json",
				"CVE-2012-0053.json",
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			r := brbundle.NewRepository(brbundle.ROption{
				OmitEnvVarFolderBundle: true,
				OmitExeBundle:          true,
				OmitEmbeddedBundle:     true,
			})
			{
				ts := initServer("./testdata/result/old-manifest")
				defer ts.Close()

				p, err := r.RegisterRemoteManifest(ts.URL, brbundle.Option{
					TempFolder:          filepath.Join(os.TempDir(), "read-test-2"),
					ResetDownloadFolder: true,
					MountPoint:          testcase.MountPoint,
				})
				if err != nil {
					assert.Nil(t, err)
					return
				}
				p.Wait()
			}

			dirs := r.Dirs()
			assert.Equal(t, testcase.Dirs, dirs)

			files := r.FilesInDir(testcase.TestDir)
			assert.Equal(t, testcase.FilesInDir, files)
		})
	}
}
