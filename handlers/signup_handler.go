package handlers

import (
	"net/http"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func SignupHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				ctx.JSON(http.StatusOK, gin.H{
					"success": true,
					"now":     time.Now(),
				})
			},
		},
	}
}
