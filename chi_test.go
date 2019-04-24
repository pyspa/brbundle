package brbundle_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brhttp"
	"github.com/stretchr/testify/assert"
)

func TestMountWithChi(t *testing.T) {
	repo := initRepo()
	r := chi.NewRouter()
	r.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, `{"status": "ok"}`)
	})
	r.NotFound(brhttp.MountFunc(brbundle.WebOption{
		Repository: repo,
	}))
	s := httptest.NewServer(r)
	defer s.Close()

	res, err := http.Get(s.URL + "/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(string(body), "hello world from root\n"))

	res2, err := http.Get(s.URL + "/api/status")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, res2.StatusCode)
}
