package brbundle

import (
	"github.com/pierrec/lz4"
	"io"
	"io/ioutil"

	"github.com/dsnet/compress/brotli"
)

const ZIPMethodLZ4 uint16 = 65535

type Decompressor func(io.Reader) io.Reader

func brotliDecompressor(input io.Reader) io.Reader {
	reader, _ := brotli.NewReader(input, nil)
	return reader
}

func nullDecompressor(input io.Reader) io.Reader {
	return input
}

func lz4Decompressor(in io.Reader) io.ReadCloser {
	return ioutil.NopCloser(lz4.NewReader(in))
}
