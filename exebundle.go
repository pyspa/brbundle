package brbundle

import (
	"os"

	"github.com/shibukawa/zipsection"
)

func NewExecutionBundle(decompressor Decompressor, decryptor Decryptor, path ...string) (FileBundle, error) {
	var filepath string
	if len(path) == 0 {
		var err error
		filepath, err = os.Executable()
		if err != nil {
			return nil, err
		}
	} else {
		filepath = path[0]
	}
	reader, closer, err := zipsection.Open(filepath)
	bundle, err := NewZipBundleFromZipReader(decompressor, decryptor, reader)
	if err != nil {
		bundle.(*ZipBundle).OnClose(func() error {
			return closer.Close()
		})
	}
	return bundle, err
}

func MustExecutionBundle(decompressor Decompressor, decryptor Decryptor, path ...string) FileBundle {
	bundle, err := NewExecutionBundle(decompressor, decryptor, path...)
	if err != nil {
		panic(err)
	}
	return bundle
}
