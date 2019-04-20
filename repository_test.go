package brbundle

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterEmbeddedBundle(t *testing.T) {
	r, err := NewRepository(ROption{
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
	r, err := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err = r.RegisterBundle("testdata/simple.pb")
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
	r, err := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err = r.RegisterBundle("testdata/simple.pb", Option{
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
	r, err := NewRepository(ROption{
		OmitEnvVarFolderBundle: true,
		OmitExeBundle:          true,
		OmitEmbeddedBundle:     true,
	})
	err = r.RegisterFolder("testdata/src-simple")
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

var bundle_628f1de9a5dbfa77bcbe37f80bc91996 = []byte(
	"PK\x03\x04\x14\x00\b\x00\x00\x00k\x03\x94N\x00\x00\x00\x00\x00\x00\x00\x00" +
		"\x00\x00\x00\x00\x05\x00\t\x00a.txtUT\x05\x00\x01\xda\xe8\xb9\\hello wo" +
		"rldPK\a\b\x85\x11J\r\v\x00\x00\x00\v\x00\x00\x00PK\x03\x04\x14\x00\b\x00" +
		"\x00\x00=\xb2\x93N\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00" +
		"\t\x00b.txtUT\x05\x00\x01\x87ʹ\\Hello World\nPK\a\b\xe3啰\f\x00\x00\x00" +
		"\f\x00\x00\x00PK\x01\x02\x14\x03\x14\x00\b\x00\x00\x00k\x03\x94N\x85\x11" +
		"J\r\v\x00\x00\x00\v\x00\x00\x00\x05\x00\t\x00\v\x00\x00\x00\x00\x00\x00" +
		"\x00\xa4\x81\x00\x00\x00\x00a.txtUT\x05\x00\x01\xda\xe8\xb9\\-b-5cb9e8d" +
		"aPK\x01\x02\x14\x03\x14\x00\b\x00\x00\x00=\xb2\x93N\xe3啰\f\x00\x00\x00" +
		"\f\x00\x00\x00\x05\x00\t\x00\v\x00\x00\x00\x00\x00\x00\x00\xa4\x81G\x00" +
		"\x00\x00b.txtUT\x05\x00\x01\x87ʹ\\-c-5cb9ca87PK\x05\x06\x00\x00\x00\x00" +
		"\x02\x00\x02\x00\x8e\x00\x00\x00\x8f\x00\x00\x00\x01\x00-")

func init() {
	RegisterEmbeddedBundle(bundle_628f1de9a5dbfa77bcbe37f80bc91996, "")
}
