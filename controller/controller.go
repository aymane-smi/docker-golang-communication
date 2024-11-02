package controller

import (
	"aymane/service"
	"context"
	"net/http"
	"strings"

	"github.com/docker/docker/client"
)

// using colsure to fix this probelm of extra params using mux
func Handler(ctx context.Context, clt *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		buf := service.DockerWriter(ctx, clt)
		output := strings.TrimSpace(buf.String())
		w.Write([]byte(output))
	}
}
