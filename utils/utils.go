package utils

import (
	"archive/tar"
	"bytes"
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
