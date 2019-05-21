package brbundle

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterEmbeddedBundle(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     false,
	})
	f, err := r.Find("a.txt")

	assert.Nil(t, err)
	assert.NotNil(t, f)
	if f != nil {
		reader, err := f.Reader()
		assert.Nil(t, err)
		b, err := ioutil.ReadAll(reader)
		assert.Nil(t, err)
		assert.Equal(t, "hello world", string(b))
	}
}

func TestRegisterPackedBundle(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err := r.RegisterBundle("testdata/simple.pb")
	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
		return
	}
	f, err := r.Find("b.txt")

	assert.Nil(t, err)
	assert.NotNil(t, f)
	if f != nil {
		assert.Equal(t, "b.txt", f.Name())
		assert.Equal(t, "/b.txt", f.Path())
	}
}

func TestRegisterPackedBundleWithMountPoint(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err := r.RegisterBundle("testdata/simple.pb", Option{
		MountPoint: "dir",
	})
	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
		return
	}
	f, err := r.Find("dir/b.txt")

	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
		return
	}
	assert.NotNil(t, f)
	if f != nil {
		assert.Equal(t, "b.txt", f.Name())
		assert.Equal(t, "/dir/b.txt", f.Path())
	}
}

func TestRepositoryFolderBundle(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err := r.RegisterFolder("testdata/src-simple")
	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
		return
	}
	f, err := r.Find("a.txt")
	assert.NotNil(t, f)
	if f != nil {
		assert.Equal(t, "a.txt", f.Name())
		assert.Equal(t, "/a.txt", f.Path())
	}
}

func TestRepositoryFolderBundle_GetLocalPath(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err := r.RegisterFolder("testdata/src-simple")
	if err != nil {
		assert.Nil(t, err)
		return
	}
	f, _ := r.Find("a.txt")
	p, err := f.GetLocalPath()
	if err != nil {
		assert.Nil(t, err)
		return
	}
	assert.True(t, strings.HasSuffix(p, f.Name()))
	content, _ := f.ReadAll()
	localContent, err := ioutil.ReadFile(p)
	if err != nil {
		assert.Nil(t, err)
		return
	}
	assert.True(t, bytes.Equal(content, localContent))
}

func TestRepositoryCache(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     false,
	})
	r.SetCacheSize(100)
	f, err := r.Find("a.txt")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	// white box test
	r.bundles[EmbeddedBundleType] = nil

	// return content from cache
	f2, err := r.Find("a.txt")
	assert.Nil(t, err)
	assert.NotNil(t, f2)
	if f2 != nil {
		reader, err := f2.Reader()
		assert.Nil(t, err)
		b, err := ioutil.ReadAll(reader)
		assert.Nil(t, err)
		assert.Equal(t, "hello world", string(b))
	}

	// now cache is empty
	r.ClearCache()
	f3, _ := r.Find("a.txt")
	assert.Nil(t, f3)
}

func TestRepositoryPackedFile_GetLocalPath(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     false,
	})
	f, _ := r.Find("a.txt")
	p, err := f.GetLocalPath()
	if err != nil {
		assert.Nil(t, err)
		return
	}
	assert.True(t, strings.HasSuffix(p, f.Name()))
	content, _ := f.ReadAll()
	localContent, err := ioutil.ReadFile(p)
	t.Log(string(content))
	t.Log(string(localContent))
	if err != nil {
		assert.Nil(t, err)
		return
	}
	assert.True(t, bytes.Equal(content, localContent))
}

func TestRegisterUnload(t *testing.T) {
	r := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	r.SetCacheSize(100)
	err := r.RegisterBundle("testdata/simple.pb")
	assert.Nil(t, err)
	if err != nil {
		t.Log(err)
		return
	}
	f, err := r.Find("b.txt")

	assert.Nil(t, err)
	assert.NotNil(t, f)

	// Unload removes cache too
	r.Unload("testdata/simple.pb")

	f2, _ := r.Find("b.txt")
	assert.Nil(t, f2)
}

var bundle_628f1de9a5dbfa77bcbe37f80bc91996 = []byte(
	"PK\x03\x04\x14\x00\b\x00\x00\x00\xb5\x1d\x94N\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x05\x00\t\x00c.txtUT\x05\x00\x01W\x17\xba\\hello w" +
		"orldPK\a\b\x85\x11J\r\v\x00\x00\x00\v\x00\x00\x00PK\x03\x04\x14\x00\b\x00" +
		"\x00\x00k\x03\x94N\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00" +
		"\t\x00a.txtUT\x05\x00\x01\xda\xe8\xb9\\hello worldPK\a\b\x85\x11J\r\v\x00" +
		"\x00\x00\v\x00\x00\x00PK\x03\x04\x14\x00\b\x00\x00\x00=\xb2\x93N\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00\t\x00b.txtUT\x05\x00\x01" +
		"\x87ʹ\\Hello World\nPK\a\b\xe3啰\f\x00\x00\x00\f\x00\x00\x00PK\x01\x02" +
		"\x14\x03\x14\x00\b\x00\x00\x00\xb5\x1d\x94N\x85\x11J\r\v\x00\x00\x00\v\x00" +
		"\x00\x00\x05\x00\t\x00\x17\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00" +
		"\x00c.txtUT\x05\x00\x01W\x17\xba\\-,b-5cba1757,text/plainPK\x01\x02\x14" +
		"\x03\x14\x00\b\x00\x00\x00k\x03\x94N\x85\x11J\r\v\x00\x00\x00\v\x00\x00" +
		"\x00\x05\x00\t\x00\x17\x00\x00\x00\x00\x00\x00\x00\xa4\x81G\x00\x00\x00" +
		"a.txtUT\x05\x00\x01\xda\xe8\xb9\\-,b-5cb9e8da,text/plainPK\x01\x02\x14\x03" +
		"\x14\x00\b\x00\x00\x00=\xb2\x93N\xe3啰\f\x00\x00\x00\f\x00\x00\x00\x05" +
		"\x00\t\x00\x17\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x8e\x00\x00\x00b.txt" +
		"UT\x05\x00\x01\x87ʹ\\-,c-5cb9ca87,text/plainPK\x05\x06\x00\x00\x00\x00" +
		"\x03\x00\x03\x00\xf9\x00\x00\x00\xd6\x00\x00\x00\x01\x00-")

func init() {
	RegisterEmbeddedBundle(bundle_628f1de9a5dbfa77bcbe37f80bc91996, "")
}

func TestPackOptions(t *testing.T) {
	testcases := []struct {
		Name        string
		BundleFile  string
		DecryptoKey string
	}{
		{
			"Snappy Compression - No Encryption",
			"sn-noe.pb",
			"",
		},
		{
			"Brotli Compression - No Encryption",
			"br-noe.pb",
			"",
		},
		{
			"Snappy Compression - AES Encryption",
			"sn-aes.pb",
			"nWKPE84p+fTc1UiMNFpPxaYFkNq44ieaNC9th8EcQC7o5c/+QRgyiKHSsc4=",
		},
		{
			"Brotli Compression - AES Encryption",
			"br-aes.pb",
			"nWKPE84p+fTc1UiMNFpPxaYFkNq44ieaNC9th8EcQC7o5c/+QRgyiKHSsc4=",
		},
	}
	testfilepaths := []string{
		"gentestdata.sh",
		"uiimage.png",
		"lena.png",
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			r := NewRepository(ROption{
				OmitEnvVarFolderBundle: true,
				OmitExeBundle:          true,
				OmitEmbeddedBundle:     true,
			})
			r.SetCacheSize(100)
			err := r.RegisterBundle(filepath.Join("testdata", testcase.BundleFile),
				Option{
					DecryptoKey: testcase.DecryptoKey,
				},
			)
			assert.Nil(t, err)
			if err != nil {
				t.Log(err)
				return
			}
			for _, testfilepath := range testfilepaths {
				t.Run(testfilepath, func(t *testing.T) {
					expected, _ := ioutil.ReadFile(filepath.Join("testdata", "src", testfilepath))
					entry, err := r.Find(testfilepath)
					assert.Nil(t, err)
					actual, err := entry.ReadAll()
					assert.Equal(t, expected, actual)
				})
			}
		})
	}
}
