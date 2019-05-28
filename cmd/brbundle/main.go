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

	generateKeyCommand = app.Command("key-gen", "Generate encryption key")

	contentFolderCommand   = app.Command("folder", "Just folder copy w/o encryption")
	contentFolderCryptoKey = contentFolderCommand.Flag("crypto", "Base64 encoded 44 bytes string to use encryption").Short('c').String()
	contentFolderTags      = contentFolderCommand.Flag("tags", "Filter files by suffix after double underscore (__)").Short('t').String()
	contentFolderDestDir   = contentFolderCommand.Arg("dest-dir", "Destination folder (folder content are removed if exists)").Required().String()
	contentFolderSourceDir = contentFolderCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	bundleCommand       = app.Command("bundle", "Append static files to an execution file")
	bundleCryptoKey     = bundleCommand.Flag("crypto", "Base64 encoded 44 bytes string to use encryption").Short('c').String()
	bundleTags          = bundleCommand.Flag("tags", "Filter files by suffix after double underscore (__)").Short('t').String()
	bundleFastCompress  = bundleCommand.Flag("fast", "Compressed by Snappy instead of Brotli").Short('f').Bool()
	bundleDirPrefix     = bundleCommand.Flag("dir-prefix", "Additional folder path added to resulting bundle contents").Short('x').String()
	bundleSpecifiedDate = bundleCommand.Flag("date", "Pseudo date of files").Short('d').String()
	bundleTargetExec    = bundleCommand.Arg("exec", "Target execution file path").Required().ExistingFile()
	bundleSourceDir     = bundleCommand.Arg("src", "Directory that contains static files").Required().ExistingDir()

	packedCommand       = app.Command("pack", "Make single packed bundle file")
	packedCryptoKey     = packedCommand.Flag("crypto", "Base64 encoded 44 bytes string to use encryption").Short('c').String()
	packedTags          = packedCommand.Flag("tags", "Filter files by suffix after double underscore (__)").Short('t').String()
	packedFastCompress  = packedCommand.Flag("fast", "Compressed by Snappy instead of Brotli").Short('f').Bool()
	packedDirPrefix     = packedCommand.Flag("dir-prefix", "Additional folder path added to resulting bundle contents").Short('x').String()
	packedSpecifiedDate = packedCommand.Flag("date", "Pseudo date of files").Short('d').String()
	packedOutputFile    = packedCommand.Arg("out-path", "Output file path").Required().OpenFile(os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	packedSourceDir     = packedCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	embeddedCommand       = app.Command("embedded", "Generate Golang code that contains")
	embeddedCryptoKey     = embeddedCommand.Flag("crypto", "Base64 encoded 44 bytes string to use encryption").Short('c').String()
	embeddedTags          = embeddedCommand.Flag("tags", "Filter files by suffix after double underscore (__)").Short('t').String()
	embeddedFastCompress  = embeddedCommand.Flag("fast", "Compressed by Snappy instead of Brotli").Short('f').Bool()
	embeddedBundleName    = embeddedCommand.Flag("name", "Bundle name to specify the bundle. It is needed to load encrypted bundle.").Short('n').String()
	packageName           = embeddedCommand.Flag("package", "Package name").Short('p').Default("main").String()
	outputFileName        = embeddedCommand.Flag("output", "Output file name").Short('o').Default("embedded-bundle.go").OpenFile(os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	embeddedDirPrefix     = embeddedCommand.Flag("dir-prefix", "Additional folder path added to resulting bundle contents").Short('x').String()
	embeddedSpecifiedDate = embeddedCommand.Flag("date", "Pseudo date of files").Short('d').String()
	embeddedSourceDir     = embeddedCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	manifestCommand       = app.Command("manifest", "Generate manifest file to sync packed bundles")
	manifestCryptoKey     = manifestCommand.Flag("crypto", "Base64 encoded 44 bytes string to use encryption").Short('c').String()
	manifestTags          = manifestCommand.Flag("tags", "Filter files by suffix after double underscore (__)").Short('t').String()
	manifestFastCompress  = manifestCommand.Flag("fast", "Compressed by Snappy instead of Brotli").Short('f').Bool()
	manifestSpecifiedDate = manifestCommand.Flag("date", "Pseudo date of files").Short('d').String()
	manifestOutputDir     = manifestCommand.Arg("out-dir", "Output directory path").Required().ExistingDir()
	manifestSourceDir     = manifestCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()
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
	kingpin.Version("1.0.0")
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
		createContentFolder(cryptoKey, *contentFolderTags, *contentFolderDestDir, *contentFolderSourceDir)
	case bundleCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*bundleCryptoKey, *bundleSpecifiedDate)
		if err != nil {
			break
		}
		appendToExec(!*bundleFastCompress, cryptoKey, *bundleTags, *bundleTargetExec, *bundleSourceDir, *bundleDirPrefix, date)
	case packedCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*packedCryptoKey, *bundleSpecifiedDate)
		if err != nil {
			break
		}
		defer (*packedOutputFile).Close()
		err = packedBundle(!*packedFastCompress, cryptoKey, *packedTags, *packedOutputFile, *packedSourceDir, *packedDirPrefix, "", date)
	case embeddedCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*embeddedCryptoKey, *embeddedSpecifiedDate)
		if err != nil {
			break
		}
		err = embedded(!*embeddedFastCompress, cryptoKey, *embeddedTags, *packageName, *outputFileName, *embeddedSourceDir, *embeddedDirPrefix, *embeddedBundleName, date)
	case manifestCommand.FullCommand():
		color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
		cryptoKey, date, err = parseKeyAndDate(*manifestCryptoKey, *manifestSpecifiedDate)
		if err != nil {
			break
		}
		err = manifest(!*manifestFastCompress, cryptoKey, *manifestTags, *manifestOutputDir, *manifestSourceDir, date)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, color.RedString("%v", err))
	}
}
