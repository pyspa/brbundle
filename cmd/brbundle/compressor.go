package main

import (
	"bytes"
	"github.com/shibukawa/brbundle"
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
}

func (c *Compressor) Compress(srcPath string, size int) (err error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

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
	c.reader, c.writer = io.Pipe()
	go func() {
		defer c.writer.Close()
		io.Copy(c.writer, src)
	}()
	cmd.Stdin = c.reader
	cmd.Stdout = c.result
	err = cmd.Run()
	c.compressedSize = c.result.Len()
	if err == nil {
		c.skipCompress = !c.shouldUseCompressedResult(size)
		if c.skipCompress {
			src, _ := os.Open(srcPath)
			defer src.Close()
			c.result.Reset()
			io.Copy(c.result, src)
		}
	}
	return
}

func (c Compressor) CompressionFlag() string {
	if !c.useLZ4 {
		return brbundle.NotToCompress
	}
	if c.skipCompress {
		return brbundle.NotToCompress
	} else {
		if c.useBrotli {
			return brbundle.UseBrotli
		} else {
			return brbundle.UseLZ4
		}
	}
}

func (c *Compressor) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, c.result)
}

func (c Compressor) shouldUseCompressedResult(uncompressed int) bool {
	compressed := c.Size()
	return (uncompressed <= 5000 && uncompressed > compressed + 100) || (uncompressed > 5000 && (int(float64(uncompressed)*0.98) > compressed))
}

func (c Compressor) Size() int {
	return c.compressedSize
}
