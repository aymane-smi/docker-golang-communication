package service

import (
	t "aymane/types"
	"aymane/utils"
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func DockerWriter(ctx context.Context, clt *client.Client, body t.Body) bytes.Buffer {
	// cases := []t.Cases{
	// 	{Input: []int{1, 2}, Expected: 3},
	// }
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

	// output := make([]byte, 0)
	// buffer := make([]byte, 1024)
	var output, stderr bytes.Buffer
	done := make(chan bool)

	go func() {
		// for {
		// 	n, err := attach.Reader.Read(buffer)
		// 	if n > 0 {
		// 		output = append(output, buffer[:n]...)
		// 	}
		// 	if err != nil {
		// 		done <- true
		// 		return
		// 	}
		// }
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

func CreateContainers() {

}
