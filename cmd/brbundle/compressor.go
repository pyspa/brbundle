package main

import (
	"archive/zip"
	"bytes"
	"github.com/pierrec/lz4"
	"io"
	"os"
	"os/exec"

	"github.com/shibukawa/brbundle"
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
		return
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
		cmd = exec.Command("lz4", "-c")
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
		if c.skipCompress || !c.useBrotli {
			src, _ := os.Open(srcPath)
			c.result.Reset()
			io.Copy(c.result, src)
		}
	}
	return
}

func (c Compressor) CompressionFlag() string {
	if !c.skipCompress && c.useBrotli {
		return brbundle.UseBrotli
	} else {
		return brbundle.NotToCompress
	}
}

func (c Compressor) ZipCompressionMethod() uint16 {
	if c.skipCompress || (!c.skipCompress && c.useBrotli) {
		return zip.Store
	} else {
		return brbundle.ZIPMethodLZ4
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

func lz4Compressor(out io.Writer) (io.WriteCloser, error) {
	return lz4.NewWriter(out), nil
}