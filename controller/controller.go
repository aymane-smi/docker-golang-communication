package controller

import (
	"aymane/service"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/client"
)

// using colsure to fix this probelm of extra params using mux
func Handler(ctx context.Context, clt *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Disposition", "inline")
		buf := service.DockerWriter(ctx, clt)
		output := strings.TrimSpace(buf.String())
		result := "the returned message is=>" + output
		fmt.Println(result)
		w.Write([]byte(result))
	}
}
