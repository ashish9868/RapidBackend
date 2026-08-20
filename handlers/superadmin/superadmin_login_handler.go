package superadmin

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/services"
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
				authService := services.NewAuthService(app)
				form, token, errors := authService.ValidateLogin(w, r, true)
				utils.Log(utils.ToJSON(errors))
				if errors != nil {
					app.RenderComponent(w, pages.LoginForm(form, errors))
					return
				}
				app.SetAuthCookie(w, token.Token, 24)
				http.Redirect(w, r, "/login?success=true", http.StatusSeeOther)
			},
		},
	}
}
