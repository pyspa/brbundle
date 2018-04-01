package main

import (
	"github.com/fatih/color"
	"io"
	"os"
	"path/filepath"
	"sync"
	"fmt"
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

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		compressor.WriteTo(encryptor)
		encryptor.Close()
		compressor.Close()
		wg.Done()
	}()

	go func() {
		callback(encryptor, etag)
		wg.Done()
	}()

	go func() {
		_, err := io.Copy(compressor, srcFile)
		if err != nil {
			color.Red("  compression error: %s\n", entry.Path, err.Error())
		}
		compressor.Close()
		srcFile.Close()
		wg.Done()
	}()

	wg.Wait()

	percent := 100
	if size != 0 {
		percent = encryptor.Size() * 100 / size
	}
	color.Green("done: %s (%d/%d=%d%%)\n", entry.Path, encryptor.Size(), size, percent)

	return nil
}
