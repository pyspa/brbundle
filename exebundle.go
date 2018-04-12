package brbundle

import (
	"os"

	"github.com/shibukawa/zipsection"
)

func NewExecutionPod(decompressor Decompressor, decryptor Decryptor, path ...string) (FilePod, error) {
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
	pod, err := NewZipPodFromZipReader(decompressor, decryptor, reader)
	if err != nil {
		pod.(*ZipPod).OnClose(func() error {
			return closer.Close()
		})
	}
	return pod, err
}

func MustExecutionPod(decompressor Decompressor, decryptor Decryptor, path ...string) FilePod {
	pod, err := NewExecutionPod(decompressor, decryptor, path...)
	if err != nil {
		panic(err)
	}
	return pod
}
