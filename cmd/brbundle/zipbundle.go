package main

import (
	"archive/zip"
	"errors"
	"io"
	"path"
	"strings"
	"sync"
	"time"
	
	"github.com/fatih/color"
)

func zipWorker(compressor *Compressor, encryptor *Encryptor, srcDirPath, dirPrefix string, date *time.Time, w *zip.Writer, lock *sync.Mutex, jobs <-chan Entry, wait chan<- struct{}) {
	for entry := range jobs {
		compressor.Init()
		encryptor.Init()

		header := &zip.FileHeader{
			Name:   cleanPath(dirPrefix, entry.Path),
			Method: zip.Store,
		}
		header.SetMode(entry.Info.Mode())
		if date != nil {
			header.Modified = *date
		} else {
			header.Modified = entry.Info.ModTime()
		}

		processInput(compressor, encryptor, srcDirPath, dirPrefix, entry, func(writerTo io.WriterTo, comment string) {
			lock.Lock()
			defer lock.Unlock()

			header.Comment = comment
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

func zipBundle(brotli bool, encryptionKey []byte, zipFile io.Writer, srcDirPath, dirPrefix, mode string, date *time.Time) error {

	_, err := NewEncryptor(encryptionKey)
	if err != nil {
		return errors.New("Can't create encryptor")
	}

	writer := zip.NewWriter(zipFile)
	var lock sync.Mutex

	color.Cyan("%s Mode (Use Brotli: %v, Use Encyrption: %v)\n\n", mode, brotli, len(encryptionKey) != 0)

	paths, _, ignored := Traverse(srcDirPath)

	wait := make(chan struct{})
	// runtime.NumCPU()
	for i := 0; i < 1; i++ {
		encryptor, _ := NewEncryptor(encryptionKey)
		go zipWorker(NewCompressor(brotli, true), encryptor, srcDirPath, dirPrefix, date, writer, &lock, paths, wait)
	}

	close(paths)

	for i := 0; i < 1; i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}
	writer.Close()
	color.Cyan("\nDone\n\n")
	return nil
}

func cleanPath(prefix, filePath string) string {
	prefix = strings.ReplaceAll(prefix, `\`, "/")
	filePath = strings.ReplaceAll(filePath, `\`, "/")
	result := path.Clean(path.Join(prefix, filePath))
	if strings.HasPrefix(result, "/") {
		return result[1:]
	} else {
		return result
	}
}
