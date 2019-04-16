package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Compressor struct {
	useBrotli      bool
	useLZ4         bool
	reader         *io.PipeReader
	writer         *io.PipeWriter
	compressed     *bytes.Buffer
	source         *os.File
	skipCompress   bool
	compressedSize int
}

func HasBrotli() bool {
	brotli := exec.Command("brotli", "--help")
	err := brotli.Run()

	return err == nil
}

func HasLZ4() bool {
	lz4 := exec.Command("lz4", "--help")
	err := lz4.Run()

	return err == nil
}

func NewCompressor(useBrotli, useLZ4 bool) *Compressor {
	return &Compressor{
		useBrotli:  useBrotli,
		useLZ4:     useLZ4,
		compressed: &bytes.Buffer{},
	}
}

func (c *Compressor) Init() {
	c.compressed.Reset()
	reader, writer := io.Pipe()
	c.reader = reader
	c.writer = writer
}

func (c *Compressor) Compress(src *os.File) (err error) {
	c.source = src
	if !c.useLZ4 {
		return
	}
	var cmd *exec.Cmd
	if c.useBrotli {
		cmd = exec.Command("brotli", "--stdout")
	} else {
		cmd = exec.Command("lz4", "-9")
	}
	go func() {
		io.Copy(c.writer, src)
		c.writer.Close()
		c.reader.Close()
	}()
	cmd.Stdin = c.reader
	cmd.Stdout = c.compressed
	err = cmd.Run()
	stat, err := src.Stat()
	c.compressedSize = c.compressed.Len()
	if err == nil {
		c.skipCompress = !c.shouldUseCompressedResult(int(stat.Size()))
	} else {
		fmt.Println(err)
	}
	return
}

func (c Compressor) CompressionFlag() string {
	if !c.useLZ4 {
		return "-"
	}
	if c.skipCompress {
		return "-"
	} else {
		if c.useBrotli {
			return "b"
		} else {
			return "l"
		}
	}
}

func (c *Compressor) WriteTo(w io.Writer) (n int64, err error) {
	defer c.source.Close()
	if c.skipCompress {
		c.source.Seek(0, 1)
		return io.Copy(w, c.source)
	} else {
		return io.Copy(w, c.compressed)
	}
}

func (c Compressor) shouldUseCompressedResult(uncompressed int) bool {
	compressed := c.Size()
	return (uncompressed-1000 > compressed) || (uncompressed > 10000 && (int(float64(uncompressed)*0.90) > compressed))
}

func (c Compressor) Size() int {
	return c.compressedSize
}
