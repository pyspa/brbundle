package brbundle_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dsnet/compress/brotli"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brhttp"
	"github.com/stretchr/testify/assert"
)

func initRepo() *brbundle.Repository {
	r := brbundle.NewRepository(brbundle.ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err := r.RegisterBundle("testdata/br-noc.pb")
	if err != nil {
		panic(err)
	}
	return r
}

func TestMount_NoBrotli(t *testing.T) {
	repo := initRepo()
	m := http.NewServeMux()
	m.Handle("/static/",
		http.StripPrefix("/static",
			brhttp.Mount(brbundle.WebOption{
				Repository: repo,
			}),
		))
	s := httptest.NewServer(m)
	defer s.Close()

	res, err := http.Get(s.URL + "/static/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(string(body), "hello world from root\n"))
}

func TestMount_Brotli(t *testing.T) {
	repo := initRepo()
	m := http.NewServeMux()
	m.Handle("/static/",
		http.StripPrefix("/static",
			brhttp.Mount(brbundle.WebOption{
				Repository: repo,
			}),
		))
	s := httptest.NewServer(m)
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
	assert.True(t, strings.HasPrefix(string(body), "hello world from root\n"))
}

func TestMountWithoutServeMux(t *testing.T) {
	repo := initRepo()
	s := httptest.NewServer(brhttp.Mount(brbundle.WebOption{
		Repository: repo,
	}))
	defer s.Close()

	res, err := http.Get(s.URL + "/rootfile.txt")
	assert.Equal(t, nil, err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(string(body), "hello world from root\n"))
}

func TestMountSPAOption(t *testing.T) {
	repo := initRepo()
	// fallback to index.html
	s := httptest.NewServer(brhttp.Mount(brbundle.WebOption{
		Repository:  repo,
		SPAFallback: "index.html",
	}))
	defer s.Close()

	res, err := http.Get(s.URL + "/static/profile")
	assert.Equal(t, nil, err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.True(t, strings.Contains(string(body), "<body>"))
}
