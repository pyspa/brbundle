package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func processInput(compressor *Compressor, encryptor *Encryptor, srcDirPath string, entry Entry, callback func(writerTo io.WriterTo, etag string)) error {
	compressor.Init()
	encryptor.Init()

	srcFile, err := os.Open(filepath.Join(srcDirPath, entry.Path))
	if err != nil {
		color.Red("read error: %s: %s\n", entry.Path, err.Error())
		return err
	}
	size := int(entry.Info.Size())
	etag := fmt.Sprintf("%x-%x", size, entry.Info.ModTime().Unix())

	err = compressor.Compress(srcFile)
	if err != nil {
		color.Red("  compression error: %s\n", entry.Path, err.Error())
	}
	srcFile.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		compressor.WriteTo(encryptor)
		encryptor.Close()
		wg.Done()
	}()

	go func() {
		comment := compressor.CompressionFlag() + encryptor.EncryptionFlag() + etag
		callback(encryptor, comment)
		wg.Done()
	}()

	wg.Wait()

	if size != 0 {
		if compressor.skipCompress {
			color.Green("done: %s (%d bytes, skip compression)\n", entry.Path, size)
		} else {
			percent := compressor.Size() * 100 / size
			color.Green("done: %s (%d bytes / %d bytes = %d%%)\n", entry.Path, compressor.Size(), size, percent)
		}
	} else {
		color.Green("done: %s (0 bytes)\n", entry.Path, compressor.Size())
	}

	return nil
}
