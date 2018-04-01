package main

import (
	"archive/zip"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"io"
	"os"
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
	size, err := detectZipOffset(file, finfo.Size())
	if err != nil {
		return err
	}
	err = file.Truncate(size)
	file.Close()
	return err
}

func detectZipOffset(file io.ReaderAt, size int64) (int64, error) {
	handlers := []func(io.ReaderAt, int64) (int64, error){
		detectZipInMacho,
		detectZipInElf,
		detectZipInPe,
	}

	for _, handler := range handlers {
		offset, err := handler(file, size)
		if err == nil {
			return offset, nil
		}
	}
	return -1, errors.New("Couldn't Open As Executable")
}

func detectZipInMacho(rda io.ReaderAt, size int64) (int64, error) {
	file, err := macho.NewFile(rda)
	if err != nil {
		return -1, err
	}

	var max int64
	for _, load := range file.Loads {
		seg, ok := load.(*macho.Segment)
		if ok {
			end := int64(seg.Offset + seg.Filesz)
			if end > max {
				max = end
			}
		}
	}

	section := io.NewSectionReader(rda, max, size-max)
	_, err = zip.NewReader(section, section.Size())
	return max, err
}

func detectZipInPe(rda io.ReaderAt, size int64) (int64, error) {
	file, err := pe.NewFile(rda)
	if err != nil {
		return -1, err
	}

	var max int64
	for _, sec := range file.Sections {
		end := int64(sec.Offset + sec.Size)
		if end > max {
			max = end
		}
	}

	section := io.NewSectionReader(rda, max, size-max)
	_, err = zip.NewReader(section, section.Size())
	return max, err
}

func detectZipInElf(rda io.ReaderAt, size int64) (int64, error) {
	file, err := elf.NewFile(rda)
	if err != nil {
		return -1, err
	}

	var max int64
	for _, sect := range file.Sections {
		if sect.Type == elf.SHT_NOBITS {
			continue
		}
		end := int64(sect.Offset + sect.Size)
		if end > max {
			max = end
		}
	}

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(rda, max, size-max)
	_, err = zip.NewReader(section, section.Size())
	return max, err
}
