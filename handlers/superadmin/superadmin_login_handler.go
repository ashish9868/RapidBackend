package superadmin

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/templates/pages"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func SuperAdminLoginHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				pages.LoginPage().Render(ctx, ctx.Writer)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				var form dto.LoginForm
				// ShouldBind checks Content-Type to select a binding engine automatically
				if err := app.BindSafely(ctx, &form); err != nil {
					app.RenderComponent(ctx, pages.LoginForm(&form, app.FormatErrors(err)))
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
					app.RenderComponent(ctx, pages.LoginForm(&form, app.FormatErrors(err)))
					return
				}

				app.RenderComponent(ctx, pages.LoginForm(&form, map[string]any{}))
			},
		},
	}
}
