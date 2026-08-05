package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/xid"
)

type ProjectForm struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Slug        string `json:"slug" form:"slug"`
}

func NewProjectsHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				projects := []models.Project{}
				count, _ := app.Bun.NewSelect().Model(&projects).Order("created_at DESC").ScanAndCount(context.Background())
				paginatedResponse(ctx, projects, count)
			},
		},
		Show: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				project := &models.Project{}
				err := app.Bun.NewSelect().Model(project).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || project.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Project not found")
					return
				}
				ctx.JSON(http.StatusOK, project)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				var form ProjectForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				form.Slug = strings.TrimSpace(form.Slug)
				if form.Slug == "" {
					form.Slug = strings.ToLower(strings.ReplaceAll(form.Title, " ", "-"))
				}
				err := validation.ValidateStruct(&form,
					validation.Field(&form.Title, validation.Required),
					validation.Field(&form.Slug, validation.Required),
				)
				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, app.FormatErrors(err))
					return
				}
				project := &models.Project{
					ID:          xid.New().String(),
					Title:       form.Title,
					Descriptiom: form.Description,
					Slug:        form.Slug,
					Settings:    map[string]any{},
				}
				_, err = app.Bun.NewInsert().Model(project).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to create project")
					return
				}
				ctx.JSON(http.StatusOK, project)
			},
		},
		Update: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				project := &models.Project{}
				err := app.Bun.NewSelect().Model(project).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || project.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Project not found")
					return
				}
				var form ProjectForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Title != "" {
					project.Title = form.Title
				}
				if form.Description != "" {
					project.Descriptiom = form.Description
				}
				if form.Slug != "" {
					project.Slug = form.Slug
				}
				_, err = app.Bun.NewUpdate().Model(project).WherePK().Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to update project")
					return
				}
				ctx.JSON(http.StatusOK, project)
			},
		},
		Delete: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				_, err := app.Bun.NewDelete().Model((*models.Project)(nil)).Where("id = ?", ctx.Param("id")).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to delete project")
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		},
	}
}
