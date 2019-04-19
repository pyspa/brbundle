package main

import (
	"github.com/fatih/color"
	"os"
	"time"
)

func appendToExec(brotli bool, ekey []byte, filePath, srcDirPath, dirPrefix string, date *time.Time) {
	stat, err := os.Stat(filePath)
	if err != nil {
		color.Red("Can't load exe file %s\n", filePath)
	}
	truncateAddedZip(filePath)
	exefile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, stat.Mode())
	defer exefile.Close()
	zipBundle(brotli, ekey, exefile, srcDirPath, dirPrefix, "Bundle to Execution", date)
}
