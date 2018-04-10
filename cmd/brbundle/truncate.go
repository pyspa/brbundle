package main

import (
	"os"

	"github.com/shibukawa/zipsection"
)

func truncateAddedZip(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	finfo, err := file.Stat()
	if err != nil {
		return err
	}
	size, err := zipsection.DetectFromReader(file, finfo.Size())
	if err != nil {
		return err
	}
	err = file.Truncate(finfo.Size() - size)
	file.Close()
	return err
}

