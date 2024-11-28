package service

import (
	t "aymane/types"
	"aymane/utils"
	"bytes"
	"context"
	"fmt"
<<<<<<< Updated upstream

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
=======
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
func CreateContainers() {

=======
// cold-start func
// check if the containers exist
// if not create one
// finally run the containers(node-php)
// it return the number of created container(it should be two [javascript , php]) + the error if existe
func CreateContainers(ctx context.Context, clt *client.Client) (int64, [2]error) {

	ArrOfLang := []string{"javascript", "php"}
	//counter is a variable that store the number of runned container after executing this method
	var counter int64 = 0
	var imageName string
	errors := [2]error{nil, nil}

	for index, lang := range ArrOfLang {
		result, err := utils.CheckExistanceOfContainer(lang, ctx, clt)
		if err != nil {
			errors[index] = err
		} else if result == false {
			//if the container doesn't existe then create a new one and started
			if lang == "javascript" {
				imageName = "node:alpine"
			} else {
				imageName = "php:8-alpine"
			}
			reader, err := clt.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				errors[index] = err
				break
			}
			io.Copy(io.Discard, reader)
			res, err := clt.ContainerCreate(ctx, &container.Config{
				Image: imageName,
				Cmd:   []string{"tail", "-f", "/dev/null"}, //using this command to prevent the container from stopping after the finish of the run
			}, &container.HostConfig{
				Resources: container.Resources{
					Memory:    int64(128 * 1024 * 1024), //for now i'm using 128 mb for memory after that we can see a better strategy
					CPUQuota:  100000,                   //using 1 cpu core in this case
					CPUPeriod: 100000,                   // for now i'm using 100% of the core it can be reduced later after a benchmark
				},
			}, nil, nil, lang)
			if err != nil {
				errors[index] = err
				break
			}
			if err := utils.StartContainer(res.ID, ctx, clt, &counter); err != nil {
				errors[index] = err
				break
			}
			//handling the case of the existance of the container
			//now check if the container is started
			//if no just started
			//the implementation it can be changed later
		} else {
			state, err := utils.CheckStateOfContainer(lang, ctx, clt)
			if err != nil {
				errors[index] = err
				break
			} else if state {
				counter++
			} else if !state {
				//if the container is stopped try to started
				//if there is any error during the launch handle the error
				if err := utils.StartContainer(lang, ctx, clt, &counter); err != nil {
					errors[index] = err
					break
				}
			}

		}
	}
	return counter, errors
>>>>>>> Stashed changes
}
