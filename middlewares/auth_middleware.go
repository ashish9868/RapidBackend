package middlewares

import (
	"strings"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/services"
	"github.com/gin-gonic/gin"
)

func NewAuthMiddleWare(app *core.App, throw401 bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authService := services.NewAuthService(app)
		token, err := ctx.Cookie(gin.AuthUserKey)
		if err != nil {
			token = ctx.GetHeader("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				app.HttpUnauthorized(ctx)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			if len(token) < 1 {
				app.HttpUnauthorized(ctx)
				return
			}
		}

		user := authService.GetUserByToken(token)
		if user != nil {
			ctx.Set(gin.AuthUserKey, user)
		}
		if user == nil && throw401 {
			app.HttpUnauthorized(ctx)
			return
		}
	}
}
