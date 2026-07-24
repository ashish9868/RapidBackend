package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func LogoutHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				app.SetAuthCookie(ctx, "", -100)
				ctx.JSON(http.StatusOK, nil)
				ctx.Abort()
			},
		},
	}
}
