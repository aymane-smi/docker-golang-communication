package utils

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
)

func CreateTar(filename string, code string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	header := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(code)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(code)); err != nil {
		return nil, err
	}

	defer tw.Close()
	return buf, nil

}

func GenerateExt(language string) (string, error) {
	if language == "php" {
		return ".php", nil
	} else if language == "javascript" {
		return ".js", nil
	} else if language == "java" {
		return ".php", nil
	} else {
		return "", errors.New("langauge not supported")
	}
}
