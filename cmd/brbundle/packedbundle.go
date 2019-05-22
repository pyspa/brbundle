package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"go.pyspa.org/brbundle"
	"io"
	"path"
	"runtime"
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
			Name:   cleanPath(dirPrefix, entry.DestPath),
			Method: compressor.ZipCompressionMethod(),
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

func packedBundle(brotli bool, encryptionKey []byte, buildTag string, outFile io.Writer, srcDirPath, dirPrefix, mode string, date *time.Time) error {

	e, err := NewEncryptor(encryptionKey)
	if err != nil {
		return errors.New("Can't create encryptor")
	}

	w := zip.NewWriter(outFile)
	w.RegisterCompressor(brbundle.ZIPMethodSnappy, snappyCompressor)
	defer w.Close()
	w.SetComment(e.EncryptionFlag())
	var lock sync.Mutex

	bt := ""
	if buildTag != "" {
		bt = fmt.Sprintf(", Build Tag: %#v", buildTag)
	}
	ec := len(encryptionKey) != 0

	if mode == "" {
		color.Cyan("Packed Bundle Mode (Use Brotli: %v, Use Encyrption: %v%s)\n\n", brotli, ec, bt)
	} else {
		color.Cyan("%s Mode (Use Brotli: %v, Use Encyrption: %v%s)\n\n", mode, brotli, ec, bt)
	}

	paths, _, ignored := Traverse(srcDirPath, buildTag)

	wait := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		encryptor, _ := NewEncryptor(encryptionKey)
		go zipWorker(NewCompressor(brotli, true), encryptor, srcDirPath, dirPrefix, date, w, &lock, paths, wait)
	}

	close(paths)

	for i := 0; i < runtime.NumCPU(); i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}
	if mode == "" {
		color.Cyan("\nDone\n\n")
	}
	return nil
}

func packedBundleShallow(brotli bool, encryptionKey []byte, buildTag string, outFile io.Writer, srcDirPath, dirPrefix, mode string, date *time.Time) error {

	e, err := NewEncryptor(encryptionKey)
	if err != nil {
		return errors.New("Can't create encryptor")
	}

	w := zip.NewWriter(outFile)
	w.RegisterCompressor(brbundle.ZIPMethodSnappy, snappyCompressor)
	defer w.Close()
	w.SetComment(e.EncryptionFlag())
	var lock sync.Mutex

	bt := ""
	if buildTag != "" {
		bt = fmt.Sprintf(", Build Tag: %#v", buildTag)
	}
	ec := len(encryptionKey) != 0

	if mode == "" {
		color.Cyan("Packed Bundle Mode (Use Brotli: %v, Use Encyrption: %v%s)\n\n", brotli, ec, bt)
	} else {
		color.Cyan("%s Mode (Use Brotli: %v, Use Encyrption: %v%s)\n\n", mode, brotli, ec, bt)
	}

	paths, _, ignored := TraverseShallow(srcDirPath, buildTag)

	wait := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		encryptor, _ := NewEncryptor(encryptionKey)
		go zipWorker(NewCompressor(brotli, true), encryptor, srcDirPath, dirPrefix, date, w, &lock, paths, wait)
	}

	close(paths)

	for i := 0; i < runtime.NumCPU(); i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}
	if mode == "" {
		color.Cyan("\nDone\n\n")
	}
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
