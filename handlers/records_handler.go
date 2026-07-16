package handlers

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func NewRecordsHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
			},
		},
	}
}
