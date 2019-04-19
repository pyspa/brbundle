package brbundle

import (
	"io"

	"github.com/dsnet/compress/brotli"
	"github.com/pierrec/lz4"
)

type Decompressor func(io.Reader) io.Reader

func brotliDecompressor(input io.Reader) io.Reader {
	reader, _ := brotli.NewReader(input, nil)
	return reader
}

func lz4Decompressor(input io.Reader) io.Reader {
	return lz4.NewReader(input)
}

func nullDecompressor(input io.Reader) io.Reader {
	return input
}
