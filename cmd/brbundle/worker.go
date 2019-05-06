package main

import (
	"fmt"
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

	path := cleanPath(dirPrefix, entry.DestPath)
	renameInfo := ""
	if entry.Path != entry.DestPath {
		renameInfo = fmt.Sprintf("Original File = %s, ", entry.Path)
	}
	sizeInfo := ""
	if size != 0 {
		if compressor.skipCompress {
			sizeInfo = fmt.Sprintf("%d bytes, skip compression", size)
		} else if compressor.Size() == 0 {
			sizeInfo = fmt.Sprintf("%d bytes", size)
		} else {
			percent := compressor.Size() * 100 / size
			sizeInfo = fmt.Sprintf("%d bytes / %d bytes = %d%%", compressor.Size(), size, percent)
		}
	} else {
		sizeInfo = "0 bytes"
	}
	color.Green("done: %s (%s%s)\n", path, renameInfo, sizeInfo)

	return nil
}
