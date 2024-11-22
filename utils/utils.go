package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func CheckExistanceOfContainer(lang string, ctx context.Context, clt *client.Client) (error, bool) {
	list, err := clt.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err, false
	}
	for _, container_ := range list {
		join_str := strings.Join(container_.Names, "")
		if join_str[1:] == lang {
			return nil, true
		}
	}
	return err, false
}

func CheckStateOfConatiner(lang string, ctx context.Context, clt *client.Client) (error, bool) {
	json, err := clt.ContainerInspect(ctx, "php")
	return err, json.State.Running
}
