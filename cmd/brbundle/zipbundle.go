package main

import (
	"archive/zip"
	"github.com/shibukawa/brbundle"
	"io"
	"os"
	"runtime"
	"github.com/fatih/color"
	"sync"
)

func zipWorker(compressor *Compressor, encryptor *Encryptor, srcDirPath string, w *zip.Writer, lock *sync.Mutex, jobs <-chan Entry, wait chan<- struct{}) {
	for entry := range jobs {
		compressor.Init()
		encryptor.Init()

		header := &zip.FileHeader{
			Name:   entry.Path,
			Method: zip.Store,
		}
		header.SetMode(entry.Info.Mode())
		header.SetModTime(entry.Info.ModTime())

		processInput(compressor, encryptor, srcDirPath, entry, func(writerTo io.WriterTo, etag string) {
			lock.Lock()
			defer lock.Unlock()

			header.Comment = etag
			f, err := w.CreateHeader(header)
			if err != nil {
				color.Red("  write file error: %s\n", entry.Path, err.Error())
			} else {
				writerTo.WriteTo(f)
			}
		})

	}
	wait <- struct{}{}
}

func zipBundle(ctype brbundle.CompressionType, etype brbundle.EncryptionType, ekey, nonce []byte, zipFile *os.File, srcDirPath string) {
	writer := zip.NewWriter(zipFile)
	var lock sync.Mutex

	color.Cyan("Zip Mode (Compression: %s, Encyrption: %s)\n\n", ctype, etype)

	paths, _, ignored := Traverse(srcDirPath)

	wait := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go zipWorker(NewCompressor(ctype), NewEncryptor(etype, ekey, nonce), srcDirPath, writer, &lock, paths, wait)
	}

	close(paths)

	for i := 0; i < runtime.NumCPU(); i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}
	writer.Close()
	color.Cyan("\nDone\n\n")
}
