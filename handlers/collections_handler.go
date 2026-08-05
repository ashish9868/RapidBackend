package handlers

import (
	"context"
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/xid"
)

type CollectionForm struct {
	ProjectID string         `json:"project_id" form:"project_id"`
	Name      string         `json:"name" form:"name"`
	Type      string         `json:"type" form:"type"`
	SortOrder int            `json:"sort_order" form:"sort_order"`
	Required  bool           `json:"required" form:"required"`
	Rules     map[string]any `json:"rules" form:"rules"`
	Options   map[string]any `json:"options" form:"options"`
}

func NewCollectionsHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				collections := []models.ProjectCollection{}
				query := app.Bun.NewSelect().Model(&collections).Order("sort_order ASC", "created_at DESC")
				if projectID := ctx.Query("project_id"); projectID != "" {
					query = query.Where("project_id = ?", projectID)
				}
				count, _ := query.ScanAndCount(context.Background())
				paginatedResponse(ctx, collections, count)
			},
		},
		Show: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				collection := &models.ProjectCollection{}
				err := app.Bun.NewSelect().Model(collection).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || collection.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Collection not found")
					return
				}
				ctx.JSON(http.StatusOK, collection)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				var form CollectionForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Type == "" {
					form.Type = "base"
				}
				err := validation.ValidateStruct(&form,
					validation.Field(&form.ProjectID, validation.Required),
					validation.Field(&form.Name, validation.Required),
				)
				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, app.FormatErrors(err))
					return
				}
				collection := &models.ProjectCollection{
					ID:        xid.New().String(),
					ProjectID: form.ProjectID,
					Name:      form.Name,
					Type:      form.Type,
					SortOrder: form.SortOrder,
					Required:  form.Required,
					Rules:     form.Rules,
					Options:   form.Options,
				}
				if collection.Rules == nil {
					collection.Rules = map[string]any{}
				}
				if collection.Options == nil {
					collection.Options = map[string]any{}
				}
				_, err = app.Bun.NewInsert().Model(collection).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to create collection")
					return
				}
				ctx.JSON(http.StatusOK, collection)
			},
		},
		Update: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				collection := &models.ProjectCollection{}
				err := app.Bun.NewSelect().Model(collection).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || collection.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Collection not found")
					return
				}
				var form CollectionForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Name != "" {
					collection.Name = form.Name
				}
				if form.Type != "" {
					collection.Type = form.Type
				}
				collection.SortOrder = form.SortOrder
				collection.Required = form.Required
				if form.Rules != nil {
					collection.Rules = form.Rules
				}
				if form.Options != nil {
					collection.Options = form.Options
				}
				_, err = app.Bun.NewUpdate().Model(collection).WherePK().Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to update collection")
					return
				}
				ctx.JSON(http.StatusOK, collection)
			},
		},
		Delete: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				id := ctx.Param("id")
				_, _ = app.Bun.NewDelete().Model((*models.ProjectCollectionRecord)(nil)).Where("collection_id = ?", id).Exec(context.Background())
				_, _ = app.Bun.NewDelete().Model((*models.ProjectCollectionField)(nil)).Where("collection_id = ?", id).Exec(context.Background())
				_, err := app.Bun.NewDelete().Model((*models.ProjectCollection)(nil)).Where("id = ?", id).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to delete collection")
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		},
	}
}
