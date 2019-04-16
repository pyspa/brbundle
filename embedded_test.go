package brbundle_test

import (
	"github.com/shibukawa/brbundle"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

var _Bundleaf286c585c49f12fc5e22571545e8442015475d2 = []byte("hello world from root\n")
var _Bundle8356acc8cb6d8db98251757dc6bfeaebeaa4c790 = []byte("hello world from subfolder\n")

// Bundle returns content bundle for brbundle FileSystem
var RawNoEncBundle = brbundle.MustEmbeddedBundle(brbundle.NullDecompressor(), brbundle.NullDecryptor(), map[string][]string{
	"/":          []string{"rootfile.txt"},
	"/subfolder": []string{"file-in-subfolder.txt"},
}, map[string]*brbundle.Entry{
	"/rootfile.txt": &brbundle.Entry{
		Path:         "/rootfile.txt",
		FileMode:     0644,
		OriginalSize: 22,
		Mtime:        time.Unix(1522526276, 1522526276000000000),
		Data:         _Bundleaf286c585c49f12fc5e22571545e8442015475d2,
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
		Data:         _Bundle8356acc8cb6d8db98251757dc6bfeaebeaa4c790,
		ETag:         "1b-5abfe816",
	},
})

func TestEmbeddedRawNoCrypto_RootFile(t *testing.T) {
	var bundle = brbundle.NewBundle(RawNoEncBundle())
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestEmbeddedRawNoCrypto_SubDirFile(t *testing.T) {
	var bundle = brbundle.NewBundle(RawNoEncBundle())
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}

var _Bundleaecb14885bc3898acc068f03a1eacdd345ad14a5 = []byte("E\xbf\xd3P)>7\fj\x91\x19\xd1B\x94\x18\x99H\x12\xe2\bZ\x11~Z\xaf>\xf7ʺ\xb9ZZ\xf6\xc8\x1e\f\xf5\xa7\xb2\xf2b\xa3\x15\xe5͋\xfc")

var _Bundlec9651f1639016275cc07ae093dea6ac9f08a835d = []byte("\fp\x8e;\x05\x9a!\x01W\xe9\t el\xdd\xf2\x00\x0e*0[\xb0s\xbdـkl\xc0e\x97\xafm?\xf4Bi`\xf4\xa8\xf9>tU\xbb\x89\x0f<>K\x9a\xfb\x14\xc6\xee")

var BrotliAESBundle = brbundle.MustEmbeddedBundle(brbundle.BrotliDecompressor(), brbundle.AESDecryptor(), map[string][]string{
	"/":          []string{"rootfile.txt"},
	"/subfolder": []string{"file-in-subfolder.txt"},
}, map[string]*brbundle.Entry{

	"/rootfile.txt": &brbundle.Entry{
		Path:         "/rootfile.txt",
		FileMode:     0644,
		OriginalSize: 22,
		Mtime:        time.Unix(1522526276, 1522526276000000000),
		Data:         _Bundleaecb14885bc3898acc068f03a1eacdd345ad14a5,
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
		Data:         _Bundlec9651f1639016275cc07ae093dea6ac9f08a835d,
		ETag:         "1b-5abfe816",
	},
})

func TestEmbeddedBrotliAES_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(BrotliAESBundle(encryptKey))
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestEmbeddedBrotliAES_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(BrotliAESBundle(encryptKey))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}

var _Bundledc000bbfe676e07dc7c16856c5b7e0ec6629f6ba = []byte("\xbe\rj\xa7\x83\xe5-\x06\xf2\a\xb0*\x88u\xa3n{\xfb\xf7\xd8\x1000\u007f\xc7F\xeb\xf1\x8b-\xbb'M\xc8\xe6\x1d\x05\xce\xf4\xfc\xa8\u05ce3\xc61KbO<\xa9\xaf\xeb\xe1\x1d\xa4\x80\x04\u007f\xed\x06")

var _Bundled1430c4ab19910e06385424a49b633072e74995a = []byte("\b2L+\xe9XF\xbf\x80\\\xad\xb2\xe9\xbf\xf5YQ:_\x17\x89\xa4\x99\xe9\x81#\xd9y\xfc\x82bs\xd8\xfe\x8f\xf7&]\xb2[\xc2\x1c\xfc\xa0\x19\x86A\xac\xccߚ\xc3G\x17j\x19\x81\xf8\xa4\x02\xb7\x0e=\xdd\xf2\xee")

var LZ4ChaChaBundle = brbundle.MustEmbeddedBundle(brbundle.LZ4Decompressor(), brbundle.ChaChaDecryptor(), map[string][]string{
	"/":          []string{"rootfile.txt"},
	"/subfolder": []string{"file-in-subfolder.txt"},
}, map[string]*brbundle.Entry{

	"/rootfile.txt": &brbundle.Entry{
		Path:         "/rootfile.txt",
		FileMode:     0644,
		OriginalSize: 22,
		Mtime:        time.Unix(1522526276, 1522526276000000000),
		Data:         _Bundledc000bbfe676e07dc7c16856c5b7e0ec6629f6ba,
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
		Data:         _Bundled1430c4ab19910e06385424a49b633072e74995a,
		ETag:         "1b-5abfe816",
	},
})

func TestEmbeddedLZ4ChaCha_RootFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(LZ4ChaChaBundle(encryptKey))
	entry, err := bundle.Find("/rootfile.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "rootfile.txt", entry.Name())
	assert.Equal(t, "/rootfile.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from root\n", string(content))
}

func TestEmbeddedLZ4ChaCha_SubDirFile(t *testing.T) {
	encryptKey := []byte("12345678123456781234567812345678")
	var bundle = brbundle.NewBundle(LZ4ChaChaBundle(encryptKey))
	entry, err := bundle.Find("/subfolder/file-in-subfolder.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "file-in-subfolder.txt", entry.Name())
	assert.Equal(t, "/subfolder/file-in-subfolder.txt", entry.Path())
	reader, err := entry.Reader()
	assert.Equal(t, nil, err)
	content, err := ioutil.ReadAll(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world from subfolder\n", string(content))
}
