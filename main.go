package main

import (
	"aymane/controller"
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	ctx := context.Background()
	clt, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer clt.Close()
	r.HandleFunc("/docker", controller.Handler(ctx, clt))
	http.ListenAndServe(":8008", r)
}
