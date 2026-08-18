package superadmin

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/templates/pages"
	"github.com/ashish9868/rapidbackend/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func SuperAdminLoginHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				utils.Log("coming here")
				app.RenderComponent(w, pages.LoginPage())
			},
		},
		Create: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				var form dto.LoginForm
				// ShouldBind checks Content-Type to select a binding engine automatically
				if err := app.BindSafely(w, r, &form); err != nil {
					app.RenderComponent(w, pages.LoginForm(&form, app.FormatErrors(err)))
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
					app.RenderComponent(w, pages.LoginForm(&form, app.FormatErrors(err)))
					return
				}

				app.RenderComponent(w, pages.LoginForm(&form, map[string]any{}))
			},
		},
	}
}
