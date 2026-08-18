package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
)

func HealthHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {

			},
		},
	}
}
