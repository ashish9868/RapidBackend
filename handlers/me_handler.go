package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func MeHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				val, exists := ctx.Get(gin.AuthUserKey)
				if !exists {
					ctx.JSON(http.StatusUnauthorized, nil)
					return
				}
				ctx.JSON(http.StatusOK, val)
				return
			},
		},
	}
}
