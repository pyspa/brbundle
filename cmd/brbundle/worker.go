package main

import (
	"io"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/shibukawa/brbundle/websupport"
)

func processInput(compressor *Compressor, encryptor *Encryptor, srcDirPath, dirPrefix string, entry Entry, callback func(writerTo io.WriterTo, etag string)) error {
	compressor.Init()
	encryptor.Init()

	srcPath := filepath.Join(srcDirPath, entry.Path)
	size := int(entry.Info.Size())

	err := compressor.Compress(srcPath, size)
	if err != nil {
		color.Red("  compression error: %s %v\n", srcPath, err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		compressor.WriteTo(encryptor)
		encryptor.Close()
		wg.Done()
	}()

	go func() {
		comment := websupport.MakeCommentString(compressor.CompressionFlag(), srcPath, entry.Info)
		callback(encryptor, comment)
		wg.Done()
	}()

	wg.Wait()

	path := cleanPath(dirPrefix, entry.Path)
	if size != 0 {
		if compressor.skipCompress {
			color.Green("done: %s (%d bytes, skip compression)\n", path, size)
		} else if compressor.Size() == 0 {
			color.Green("done: %s (%d bytes)\n", path, size)
		} else {
			percent := compressor.Size() * 100 / size
			color.Green("done: %s (%d bytes / %d bytes = %d%%)\n", path, compressor.Size(), size, percent)
		}
	} else {
		color.Green("done: %s (0 bytes)\n", path, compressor.Size())
	}

	return nil
}
