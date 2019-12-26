package util

import (
	"bytes"
	"io"
	"os"
)

func ReadBytes(fp io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, fp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadFileFromPath(path string) ([]byte, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return ReadBytes(fp)
}
