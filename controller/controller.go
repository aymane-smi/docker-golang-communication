package controller

import (
	"aymane/service"
	t "aymane/types"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/docker/client"
)

// using colsure to fix this probelm of extra params using mux
func Handler(ctx context.Context, clt *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body t.Body
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			panic(err.Error())
		}
		buf := service.DockerWriter(ctx, clt, body.Language)
		output := strings.TrimSpace(buf.String())
		w.Write([]byte(output))
	}
}
