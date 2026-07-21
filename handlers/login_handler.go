package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/services"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				time.Sleep(5 * time.Second)
				var form LoginForm
				// ShouldBind checks Content-Type to select a binding engine automatically
				if err := ctx.ShouldBindJSON(&form); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"errors": err,
					})
					return
				}

				err := validation.ValidateStruct(&form,
					validation.Field(&form.Email,
						validation.Required.Error("Username is required"),
						is.Email.Error("Please Provide a valid Email"),
					),
					validation.Field(&form.Password,
						validation.Required.Error("Password is required"),
					),
				)

				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, err)
					return
				}
				authService := services.NewAuthService(app)
				token := authService.LoginByEmail(form.Email, form.Password)

				fmt.Println(token)

				if token == nil {
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{
						"global": "Unable to login, credentials are invalid.",
					})
					return
				}
				ctx.SetCookie("gin_cookie", token.Token, 3600, "/", "localhost", false, true)
				ctx.JSON(http.StatusOK, token)
			},
		},
	}
}
