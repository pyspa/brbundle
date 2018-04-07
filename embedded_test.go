package brbundle_test

import (
	"testing"
	"github.com/ToQoz/gopwt/assert"
	"github.com/shibukawa/brbundle"
	"time"
	"io/ioutil"
)

var _Podaf286c585c49f12fc5e22571545e8442015475d2 = []byte("hello world from root\n")
var _Pod8356acc8cb6d8db98251757dc6bfeaebeaa4c790 = []byte("hello world from subfolder\n")

// Pod returns content pod for brbundle FileSystem
var RawNoEncPod = brbundle.MustEmbeddedPod(brbundle.NullDecompressor(), brbundle.NullDecryptor(), map[string][]string{
	"/":          []string{"rootfile.txt"},
	"/subfolder": []string{"file-in-subfolder.txt"},
}, map[string]*brbundle.Entry{
	"/rootfile.txt": &brbundle.Entry{
		Path:         "/rootfile.txt",
		FileMode:     0644,
		OriginalSize: 22,
		Mtime:        time.Unix(1522526276, 1522526276000000000),
		Data:         _Podaf286c585c49f12fc5e22571545e8442015475d2,
		ETag:         "16-5abfe844",
	},

	"/subfolder": &brbundle.Entry{
		Path:         "/subfolder",
		FileMode:     020000000755,
		OriginalSize: 136,
		Mtime:        time.Unix(1522526252, 1522526252000000000),
		Data:         nil,
		ETag:         "",
	},

	"/subfolder/file-in-subfolder.txt": &brbundle.Entry{
		Path:         "/subfolder/file-in-subfolder.txt",
		FileMode:     0644,
		OriginalSize: 27,
		Mtime:        time.Unix(1522526230, 1522526230000000000),
		Data:         _Pod8356acc8cb6d8db98251757dc6bfeaebeaa4c790,
		ETag:         "1b-5abfe816",
	},
})


func TestEmbeddedRawNoCrypto_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(RawNoEncPod())
	entry, err := bundle.Find("/rootfile.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "rootfile.txt")
	reader, err := entry.RawReader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from root\n")
}

func TestEmbeddedRawNoCrypto_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(RawNoEncPod())
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.OK(t, err == nil)
	assert.OK(t, entry.Name() == "file-in-subfolder.txt")
	reader, err := entry.RawReader()
	assert.OK(t, err == nil)
	content, err := ioutil.ReadAll(reader)
	assert.OK(t, err == nil)
	assert.OK(t, string(content) == "hello world from subfolder\n")
}
