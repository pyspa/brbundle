package main

import (
	"io"
	"os/exec"

	//"github.com/pierrec/lz4"
	"github.com/shibukawa/brbundle"
)

type Compressor struct {
	ctype  brbundle.CompressionType
	reader *io.PipeReader
	writer *io.PipeWriter
}

func HasBrotli() bool {
	brotli := exec.Command("brotli", "--help")
	err := brotli.Run()

	return err == nil
}

func NewCompressor(ctype brbundle.CompressionType) *Compressor {
	return &Compressor{
		ctype,
		nil,
		nil,
	}
}

func (c *Compressor) Init() {
	reader, writer := io.Pipe()
	c.reader = reader
	c.writer = writer
}

func (c *Compressor) Write(data []byte) (n int, err error) {
	n, err = c.writer.Write(data)
	return
}

func (c *Compressor) Close() {
	c.writer.Close()
}

func (c *Compressor) WriteTo(w io.Writer) (n int64, err error) {
	switch c.ctype {
	case brbundle.NoCompression:
		n, err = io.Copy(w, c.reader)
	case brbundle.Brotli:
		brotli := exec.Command("brotli", "--stdout")
		brotli.Stdin = c.reader
		brotli.Stdout = w
		err = brotli.Run()
	case brbundle.LZ4:
		brotli := exec.Command("lz4", "-9")
		brotli.Stdin = c.reader
		brotli.Stdout = w
		err = brotli.Run()
	}
	return
}
