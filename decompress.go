package brbundle

import (
	"github.com/dsnet/compress/brotli"
	"io"
)

type Decompressor interface {
	Decompress(input io.Reader) io.Reader
}

type brotliDecompressor struct{}

func (b *brotliDecompressor) Decompress(input io.Reader) io.Reader {
	reader, _ := brotli.NewReader(input, nil)
	return reader
}

func BrotliDecompressor() Decompressor {
	return &brotliDecompressor{}
}

type nullDecompressor struct{}

func (b *nullDecompressor) Decompress(input io.Reader) io.Reader {
	return input
}

func NullDecompressor() Decompressor {
	return &nullDecompressor{}
}
