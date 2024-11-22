package main

import (
	"aymane/controller"
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	ctx := context.Background()
	clt, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer clt.Close()
	r.HandleFunc("/docker", controller.Handler(ctx, clt)).Methods("POST")
	handler := cors.Default().Handler(r)
	fmt.Println("server started at port 8008")
	http.ListenAndServe(":8008", handler)
}
