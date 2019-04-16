package main

import (
	"fmt"
	"io"
)

type ReaderMonitor struct {
	reader io.Reader
	label  string
}

func (r *ReaderMonitor) Read(p []byte) (n int, err error) {
	fmt.Printf("[%s] start read\n", r.label)
	n, err = r.reader.Read(p)
	fmt.Printf("[%s] end read %d bytes\n", r.label, n)
	return n, err
}

func NewReaderMonitor(reader io.Reader, label string) io.Reader {
	return &ReaderMonitor{
		reader: reader,
		label:  label,
	}
}

type WriterMonitor struct {
	writer io.Writer
	label  string
}

func (r *WriterMonitor) Write(p []byte) (n int, err error) {
	fmt.Printf("[%s] start write\n", r.label)
	n, err = r.writer.Write(p)
	fmt.Printf("[%s] end write %d bytes\n", r.label, n)
	return n, err
}

func NewWriterMonitor(writer io.Writer, label string) io.Writer {
	return &WriterMonitor{
		writer: writer,
		label:  label,
	}
}
