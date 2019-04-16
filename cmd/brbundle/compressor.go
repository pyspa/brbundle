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
	result         *bytes.Buffer
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
		useBrotli: useBrotli,
		useLZ4:    useLZ4,
		result:    &bytes.Buffer{},
	}
}

func (c *Compressor) Init() {
	c.result.Reset()
	reader, writer := io.Pipe()
	c.reader = reader
	c.writer = writer
}

func (c *Compressor) Compress(src *os.File) (err error) {
	if !c.useLZ4 {
		io.Copy(c.result, src)
		return
	}
	var cmd *exec.Cmd
	if c.useBrotli {
		cmd = exec.Command("brotli", "--stdout")
	} else {
		cmd = exec.Command("lz4", "-9", "-c")
	}
	go func() {
		io.Copy(c.writer, src)
		c.writer.Close()
		c.reader.Close()
	}()
	cmd.Stdin = c.reader
	cmd.Stdout = c.result
	err = cmd.Run()
	stat, err := src.Stat()
	c.compressedSize = c.result.Len()
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
	return io.Copy(w, c.result)
}

func (c Compressor) shouldUseCompressedResult(uncompressed int) bool {
	compressed := c.Size()
	return (uncompressed-1000 > compressed) || (uncompressed > 10000 && (int(float64(uncompressed)*0.90) > compressed))
}

func (c Compressor) Size() int {
	return c.compressedSize
}
