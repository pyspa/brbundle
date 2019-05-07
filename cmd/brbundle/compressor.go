package main

import (
	"archive/zip"
	"bytes"
	"github.com/golang/snappy"
	"io"
	"os"
	"os/exec"

	"github.com/shibukawa/brbundle"
)

type Compressor struct {
	useBrotli      bool
	useSnappy      bool
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

func NewCompressor(useBrotli, useSnappy bool) *Compressor {
	return &Compressor{
		useBrotli: useBrotli,
		useSnappy: useSnappy,
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

	if !c.useSnappy {
		io.Copy(c.result, src)
		return
	} else if c.useBrotli {
		cmd := exec.Command("brotli", "--stdout")
		c.reader, c.writer = io.Pipe()
		go func() {
			defer c.writer.Close()
			io.Copy(c.writer, src)
		}()
		cmd.Stdin = c.reader
		cmd.Stdout = c.result
		err = cmd.Run()
		c.compressedSize = c.result.Len()
	} else {
		w := snappy.NewBufferedWriter(c.result)
		var size int64
		size, err = io.Copy(w, src)
		c.compressedSize = int(size)
	}
	if err == nil {
		c.skipCompress = !shouldUseCompressedResult(c.compressedSize, size)
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
		return brbundle.ZIPMethodSnappy
	}
}

func (c *Compressor) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, c.result)
}

func shouldUseCompressedResult(compressed, uncompressed int) bool {
	return (uncompressed <= 5000 && uncompressed > compressed + 100) || (uncompressed > 5000 && (int(float64(uncompressed)*0.98) > compressed))
}

func (c Compressor) Size() int {
	return c.compressedSize
}

func snappyCompressor(out io.Writer) (io.WriteCloser, error) {
	return snappy.NewBufferedWriter(out), nil
}