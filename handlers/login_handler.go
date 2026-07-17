package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func LoginHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				data := map[string]any{}
				ctx.HTML(http.StatusOK, "public.contact", data)
			},
		},
	}
}
