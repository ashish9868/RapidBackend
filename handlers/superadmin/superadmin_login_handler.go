package superadmin

import (
	"net/http"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/services"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/templates/pages"
	"github.com/ashish9868/rapidbackend/utils"
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
				time.Sleep(5 * time.Second)
				form := &dto.LoginForm{}
				if err := app.BindSafely(w, r, form); err != nil {
					app.RenderComponent(w, pages.LoginForm(form, err))
					return
				}
				authService := services.NewAuthService(app)
				token, errors := authService.ValidateLogin(form, true)
				utils.Log(utils.ToJSON(errors))
				if errors != nil {
					app.RenderComponent(w, pages.LoginForm(form, errors))
					return
				}
				app.SetAuthCookie(w, token.Token, 24)
				app.Redirect(w, r, "/dashboard")
			},
		},
	}
}
