package brbundle

import (
	"io"
	"io/ioutil"

	"github.com/dsnet/compress/brotli"
	"github.com/golang/snappy"
)

const ZIPMethodSnappy uint16 = 65535

type Decompressor func(io.Reader) io.Reader

func brotliDecompressor(input io.Reader) io.Reader {
	reader, _ := brotli.NewReader(input, nil)
	return reader
}

func nullDecompressor(input io.Reader) io.Reader {
	return input
}

func snappyDecompressor(in io.Reader) io.ReadCloser {
	return ioutil.NopCloser(snappy.NewReader(in))
}
