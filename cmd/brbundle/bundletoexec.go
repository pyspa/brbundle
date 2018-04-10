package main

import (
	"github.com/shibukawa/brbundle"
	"github.com/fatih/color"
	"os"
)

func appendToExec(ctype brbundle.CompressionType, etype brbundle.EncryptionType, ekey, nonce []byte, filePath, srcDirPath string) {
	stat, err := os.Stat(filePath)
	if err != nil {
		color.Red("Can't load exe file %s\n", filePath)
	}
	truncateAddedZip(filePath)
	exefile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, stat.Mode())
	zipBundle(ctype, etype, ekey, nonce, exefile, srcDirPath, "Bundle to Execution")
}
