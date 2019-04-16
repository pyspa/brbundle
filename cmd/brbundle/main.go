package main

import (
	"fmt"
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("brbundle", "Static file bundle tool with high compression ratio and encryption")

	generateKeyCommand = app.Command("generate-key", "Generate encryption key")

	contentFolderCommand   = app.Command("content", "Just folder copy encryption")
	contentFolderCryptoKey = contentFolderCommand.Flag("crypto", "base64 encoded 44 bytes string to use encryption").Short('c').String()
	contentFolderDestDir   = contentFolderCommand.Arg("dest-dir", "Destination folder (folder content are removed if exists)").Required().String()
	contentFolderSourceDir = contentFolderCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	bundleCommand       = app.Command("bundle", "Append static files to an execution file")
	bundleCryptoKey     = bundleCommand.Flag("crypto", "base64 encoded 44 bytes string to use encryption").Short('c').String()
	bundleCompress      = bundleCommand.Flag("compress", "Compressed by Brotli").Short('z').Bool()
	bundleSpecifiedDate = bundleCommand.Flag("date", "Pseudo date of files").Short('d').String()
	bundleTargetExec    = bundleCommand.Arg("exec", "Target execution file path").Required().ExistingFile()
	bundleSourceDir     = bundleCommand.Arg("src", "Directory that contains static files").Required().ExistingDir()

	zipCommand       = app.Command("zip", "Make single zip file")
	zipCryptoKey     = zipCommand.Flag("crypto", "base64 encoded 44 bytes string to use encryption").Short('c').String()
	zipCompress      = zipCommand.Flag("compress", "Compressed by Brotli").Short('z').Bool()
	zipSpecifiedDate = zipCommand.Flag("date", "Pseudo date of files").Short('d').String()
	zipOutputFile    = zipCommand.Arg("zip-path", "Output zip file path").Required().OpenFile(os.O_CREATE|os.O_WRONLY, 0644)
	zipSourceDir     = zipCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	embeddedCommand       = app.Command("embedded", "Generate Golang code that contains")
	embeddedCryptoKey     = embeddedCommand.Flag("crypto", "base64 encoded 44 bytes string to use encryption").Short('c').String()
	embeddedCompress      = embeddedCommand.Flag("compress", "Compressed by Brotli").Short('z').Bool()
	packageName           = embeddedCommand.Flag("package", "Package name").Short('p').Default("main").String()
	outputFileName        = embeddedCommand.Flag("output", "Output file name").Short('o').Default("embedded-bundle.go").OpenFile(os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	embeddedSpecifiedDate = embeddedCommand.Flag("date", "Pseudo date of files").Short('d').String()
	embeddedSourceDir     = embeddedCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()
)

func parseKeyAndDate(keySrc, dateSrc string) (cryptoKey []byte, date *time.Time, err error) {
	cryptoKey, err = decodeEncryptKey(keySrc)
	if err != nil {
		return
	}
	if dateSrc == "" {
		return
	}
	var parsedDate time.Time
	parsedDate, err = dateparse.ParseStrict(dateSrc)
	if err != nil {
		return
	}
	date = &parsedDate
	return
}

func main() {
	kingpin.Version("0.0.2")
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	var err error
	var cryptoKey []byte
	var date *time.Time
	switch parse {
	case generateKeyCommand.FullCommand():
		generateKey()
	case contentFolderCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*contentFolderCryptoKey, "")
		if err != nil {
			break
		}
		createContentFolder(cryptoKey, *contentFolderDestDir, *contentFolderSourceDir)
	case bundleCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*bundleCryptoKey, *bundleSpecifiedDate)
		if err != nil {
			break
		}
		appendToExec(*bundleCompress, cryptoKey, *bundleTargetExec, *bundleSourceDir, date)
	case zipCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, err = decodeEncryptKey(*zipCryptoKey)
		if err != nil {
			break
		}
		zipBundle(*zipCompress, cryptoKey, *zipOutputFile, *zipSourceDir, "Zip", date)
	case embeddedCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*embeddedCryptoKey, *embeddedSpecifiedDate)
		if err != nil {
			break
		}
		embedded(*embeddedCompress, cryptoKey, *packageName, *outputFileName, *embeddedSourceDir, date)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, color.RedString("%v", err))
	}
}
