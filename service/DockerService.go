package service

import (
	t "aymane/types"
	"aymane/utils"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func DockerWriter(ctx context.Context, clt *client.Client, body t.Body) bytes.Buffer {
	ext, err := utils.GenerateExt(body.Language)
	if err != nil {
		panic(err.Error())
	}
	tarFile, err := utils.CreateTar("test"+ext, utils.InitTemplatePhp(body.Code, body.Name, body.TestCases))
	if err := clt.CopyToContainer(ctx, body.Language, "/tmp", tarFile, types.CopyToContainerOptions{}); err != nil {
		panic(err)
	}
	exec, err := clt.ContainerExecCreate(ctx, body.Language, types.ExecConfig{
		// save this command later for the testing "node", "-e", "console.log('hello', typeof 1)"
		Cmd:          []string{body.Language, "/tmp/test" + ext},
		AttachStdin:  true,
		AttachStdout: true,
	})

	if err != nil {
		panic(err)
	}

	attach, err := clt.ContainerExecAttach(ctx, exec.ID, container.ExecStartOptions{})
	if err != nil {
		panic(err)
	}
	defer attach.Close()

	var output, stderr bytes.Buffer
	done := make(chan bool)

	go func() {
		stdcopy.StdCopy(&output, &stderr, attach.Reader)
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("execution completed")
	case <-ctx.Done():
		panic("execution terminated due to timeout")
	}

	inspect, err := clt.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		panic(err)
	}
	if inspect.ExitCode != 0 {
		fmt.Println(err.Error())
		panic("error during the execution")
	}
	return output
}

// cold-start func
// check if the containers exist
// if not create one
// finally run the containers(node-php)
// it return the number of created container(it should be two [javascript , php]) + the error if existe
func CreateContainers(ctx context.Context, clt *client.Client) (int64, error) {
	ArrOfLang := []string{"javascript", "php"}
	var counter int64 = 0
	var imageName string
	//i should handle the case of the err where it should be a global variable and return it only at the end of the function
	for _, lang := range ArrOfLang {
		result, err := utils.CheckExistanceOfContainer(lang, ctx, clt)
		if err != nil {
			return counter, err
		} else if result == false {
			//the container doesn't existe then create a new one and started
			if lang == "javascript" {
				imageName = "node:alpine"
			} else {
				imageName = "php:8-alpine"
			}
			reader, err := clt.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				return counter, err
			}
			io.Copy(io.Discard, reader)
			res, err := clt.ContainerCreate(ctx, &container.Config{
				Image: imageName,
			}, &container.HostConfig{
				Resources: container.Resources{
					Memory:    int64(128 * 1024 * 1024), //for now i'm using 128 mb for memory after that we can see a better strategy
					CPUQuota:  100000,                   //using 1 cpu core in this case
					CPUPeriod: 100000,                   // for now i'm using 100% of the core it can be reduced later after a benchmark
				},
			}, nil, nil, lang)
			if err != nil {
				return counter, err
			}
			if err := utils.StartContainer(res.ID, ctx, clt, &counter); err != nil {
				return counter, err
			}
			//handling the case of the existance of the container
			//now check if the container is started
			//if now just started
		} else {

		}
	}
	return counter, nil
}
