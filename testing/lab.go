package main

import (
	"aymane/service"
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	clt, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	fmt.Println(service.CreateContainers(ctx, clt))
}
