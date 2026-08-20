package handlers

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/templates/pages"
)

func DashboardHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				app.RenderComponent(w, pages.DashboardHomePage(app.GetUserFromRequest(r)))
			},
		},
	}
}
