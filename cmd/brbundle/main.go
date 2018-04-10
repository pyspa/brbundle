package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/shibukawa/brbundle"
	"gopkg.in/alecthomas/kingpin.v2"
)

// https://deeeet.com/writing/2015/11/10/go-crypto/

var (
	app = kingpin.New("brbundle", "Static file bundle tool with high compression ratio and encryption")

	cryptoKey  = app.Flag("crypto-key", "32byte string to use encryption").Short('k').String()
	cryptoType = app.Flag("cipher", "Encryption algorithm").Short('c').Enum("AES", "chacha")

	compressorType = app.Flag("compressor", "Compressor type").Short('z').Default("br").Enum("br", "lz4", "raw")

	bundleCommand    = app.Command("bundle", "Append static files to an execution file")
	bundleTargetExec = bundleCommand.Arg("exec", "Target execution file path").Required().ExistingFile()
	bundleSourceDir  = bundleCommand.Arg("src", "Directory that contains static files").Required().ExistingDir()

	zipCommand    = app.Command("zip", "Make single zip file")
	zipOutputFile = zipCommand.Arg("zip-path", "Output zip file path").Required().OpenFile(os.O_CREATE|os.O_WRONLY, 0644)
	zipSourceDir  = zipCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	contentFolderCommand   = app.Command("content", "Just compress and/or encryption folder")
	contentFolderDestDir   = contentFolderCommand.Arg("dest-dir", "Destination folder (folder content are removed if exists)").Required().String()
	contentFolderSourceDir = contentFolderCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()

	embeddedCommand   = app.Command("embedded", "Generate Golang code that contains")
	packageName       = embeddedCommand.Flag("package", "Package name").Short('p').Default("main").String()
	variableName      = embeddedCommand.Flag("variable", "Variable name").Short('v').Default("Pod").String()
	outputFileName    = embeddedCommand.Flag("output", "Output file name").Short('o').Default("embedded-bundle.go").OpenFile(os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	embeddedSourceDir = embeddedCommand.Arg("src-dir", "Directory that contains static files").Required().ExistingDir()
)

func main() {
	kingpin.Version("0.0.1")
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	// warning! it is for test only
	nonce, _ := os.LookupEnv("STATIC_NONCE_FOR_TEST")

	if *cryptoKey == "" && *cryptoType != "" {
		app.Fatalf("Compressor '%s' is specified, but crypto-key is missing", *compressorType)
	}
	cryptoKeyBytes := []byte(*cryptoKey)
	if len(cryptoKeyBytes) != 0 && len(cryptoKeyBytes) != 32 {
		app.Fatalf("crypto-key should be 32byte string, but %d", len(*cryptoKey))
	}

	ctype := brbundle.NoCompression

	if contentFolderCommand.FullCommand() == "" {
		switch *compressorType {
		case "br":
			ctype = brbundle.Brotli
			if !HasBrotli() {
				color.Red("Can't run brotli on this environtent")
				os.Exit(1)
			}
		case "lz4":
			ctype = brbundle.LZ4
		default:
			ctype = brbundle.NoCompression
		}
	}

	var etype brbundle.EncryptionType

	switch *cryptoType {
	case "AES":
		etype = brbundle.AES
	case "chacha":
		etype = brbundle.ChaCha20Poly1305
	default:
		etype = brbundle.NoEncryption
	}

	color.HiBlue("\nbrbundle by Yoshiki Shibukawa\n\n")
	switch parse {
	case bundleCommand.FullCommand():
		appendToExec(ctype, etype, []byte(*cryptoKey), []byte(nonce), *bundleTargetExec, *bundleSourceDir)
	case zipCommand.FullCommand():
		zipBundle(ctype, etype, []byte(*cryptoKey), []byte(nonce), *zipOutputFile, *zipSourceDir, "Zip")
	case contentFolderCommand.FullCommand():
		createContentFolder(etype, []byte(*cryptoKey), []byte(nonce), *contentFolderDestDir, *contentFolderSourceDir)
	case embeddedCommand.FullCommand():
		embedded(ctype, etype, []byte(*cryptoKey), []byte(nonce), *packageName, *variableName, *outputFileName, *embeddedSourceDir)
	}
}
