package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/golang/snappy"
	"go.pyspa.org/brbundle"
)

type Compressor struct {
	useBrotli      bool
	useSnappy      bool
	result         *bytes.Buffer
	skipCompress   bool
	compressedSize int
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
		err = compressBrotli(c.result, src)
		c.compressedSize = c.result.Len()
	} else {
		w := snappy.NewBufferedWriter(c.result)
		_, err = io.Copy(w, src)
		w.Close()
		c.compressedSize = c.result.Len()
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
	return (uncompressed <= 5000 && uncompressed > compressed+100) || (uncompressed > 5000 && (int(float64(uncompressed)*0.98) > compressed))
}

func (c Compressor) Size() int {
	return c.compressedSize
}

func snappyCompressor(out io.Writer) (io.WriteCloser, error) {
	return snappy.NewBufferedWriter(out), nil
}

func compressBrotli(result io.Writer, src io.Reader) error {
	cmd := exec.Command("brotli", "--stdout")
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		io.Copy(writer, src)
	}()
	cmd.Stdin = reader
	cmd.Stdout = result
	return cmd.Run()
}
