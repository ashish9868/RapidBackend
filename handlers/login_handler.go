package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
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
						is.Email,
					),
					validation.Field(&form.Password,
						validation.Required.Error("Password is required"),
						validation.Length(8, 100),
					),
				)

				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, app.ErrorJson(form, err))
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"status": "logged in"})
			},
		},
	}
}
