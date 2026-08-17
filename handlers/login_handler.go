package handlers

import (
	"errors"
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func LoginHandler(collection string) *core.ResourceHandler {
	return &core.ResourceHandler{
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				view := "fragment.login"
				var form dto.LoginForm
				// ShouldBind checks Content-Type to select a binding engine automatically
				if err := app.BindSafely(ctx, &form); err != nil {
					app.SendResponse(ctx, core.Response{
						View:  view,
						Code:  http.StatusBadRequest,
						Error: err,
					})
					return
				}

				err := validation.ValidateStruct(&form,
					validation.Field(&form.Email,
						validation.Required.Error("Email is required"),
						is.Email.Error("Please Provide a valid Email"),
					),
					validation.Field(&form.Password,
						validation.Required.Error("Password is required"),
					),
				)

				if err != nil {
					app.SendResponse(ctx, core.Response{
						View:     view,
						Code:     http.StatusUnprocessableEntity,
						FormData: form,
						Error:    err,
					})
					return
				}
				user := app.AuthService.LoginByEmail(form.Email, form.Password, collection)
				if user == nil {
					app.SendResponse(ctx, core.Response{
						View:     view,
						Code:     http.StatusUnprocessableEntity,
						FormData: form,
						Error:    errors.New("Unable to login, credentials are invalid."),
					})
					return
				}
				app.SetAuthCookie(ctx, *user.Token, 3600)
				if app.IsHTMX(ctx) {
					app.SendResponse(ctx, core.Response{
						View: view,
						Code: http.StatusOK,
						Data: user,
					})
					return
				}
				ctx.JSON(http.StatusOK, user)
			},
		},
	}
}
