package brbundle

import (
	"os"

	"github.com/shibukawa/zipsection"
)

func NewExecutionPod(decompressor Decompressor, decryptor Decryptor) (FilePod, error) {
	filepath, err := os.Executable()
	if err != nil {
		return nil, err
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

func MustExecutionPod(decompressor Decompressor, decryptor Decryptor) FilePod {
	pod, err := NewExecutionPod(decompressor, decryptor)
	if err != nil {
		panic(err)
	}
	return pod
}
