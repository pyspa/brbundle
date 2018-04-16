package brbundle

import (
	"github.com/dsnet/compress/brotli"
	"github.com/pierrec/lz4"
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

type lz4Decompressor struct{}

func (l *lz4Decompressor) Decompress(input io.Reader) io.Reader {
	return lz4.NewReader(input)
}

func LZ4Decompressor() Decompressor {
	return &lz4Decompressor{}
}

type nullDecompressor struct{}

func (b *nullDecompressor) Decompress(input io.Reader) io.Reader {
	return input
}

func NullDecompressor() Decompressor {
	return &nullDecompressor{}
}
