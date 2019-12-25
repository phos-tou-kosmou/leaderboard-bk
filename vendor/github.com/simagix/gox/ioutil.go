// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

// NewFileReader returns a reader from either a gzip or plain file
func NewFileReader(filename string) (*bufio.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewReader(file)
}

// NewReader returns a reader from either a gzip or plain file
func NewReader(file *os.File) (*bufio.Reader, error) {
	var buf []byte
	var err error
	var reader *bufio.Reader

	reader = bufio.NewReader(file)
	if buf, err = reader.Peek(2); err != nil {
		return reader, err
	}
	file.Seek(0, 0)

	if buf[0] == 31 && buf[1] == 139 {
		var zreader *gzip.Reader
		if zreader, err = gzip.NewReader(file); err != nil {
			return reader, err
		}
		reader = bufio.NewReader(zreader)
	} else {
		reader = bufio.NewReader(file)
	}

	return reader, nil
}

// CountLines count number of '\n'
func CountLines(reader *bufio.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}
	lineCounts := 0
	for {
		c, err := reader.Read(buf)
		lineCounts += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return lineCounts, err

		case err != nil:
			return lineCounts, err
		}
	}
}

// OutputGzipped writes doc to a gzipped file
func OutputGzipped(b []byte, filename string) error {
	var zbuf bytes.Buffer
	var n int
	var err error
	gz := gzip.NewWriter(&zbuf)
	nw := 0
	for nw < len(b) {
		if n, err = gz.Write(b); err != nil {
			return err
		}
		nw += n
	}
	gz.Close() // flushing the bytes to the buffer.
	return ioutil.WriteFile(filename, zbuf.Bytes(), 0644)
}

// ReadAll reads from a file and return bytes
func ReadAll(file *os.File) ([]byte, error) {
	var err error
	var b []byte
	var reader *bufio.Reader

	if reader, err = NewReader(file); err != nil {
		return b, err
	}

	return ioutil.ReadAll(reader)
}
