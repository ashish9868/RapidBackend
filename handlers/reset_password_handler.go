package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ResetPasswordForm struct {
	Email string `json:"email"`
}

func ResetPasswordHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				var form ResetPasswordForm
				if err := ctx.ShouldBindJSON(&form); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"errors": err,
					})
					return
				}

				err := validation.ValidateStruct(&form,
					validation.Field(&form.Email,
						validation.Required.Error("Email is required"),
						is.Email.Error("Please Provide a valid Email"),
					),
				)

				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, err)
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})

			},
		},
	}
}
