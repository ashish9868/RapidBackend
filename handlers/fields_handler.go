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

type FieldForm struct {
	Name         string         `json:"name" form:"name"`
	Type         string         `json:"type" form:"type"`
	IsRequired   bool           `json:"is_required" form:"is_required"`
	IsIndexed    bool           `json:"is_indexed" form:"is_indexed"`
	IsUnique     bool           `json:"is_unique" form:"is_unique"`
	IsSortable   bool           `json:"is_sortable" form:"is_sortable"`
	IsFilterable bool           `json:"is_filterable" form:"is_filterable"`
	Options      map[string]any `json:"options" form:"options"`
}

func NewFieldsHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				collectionID := ctx.Param("collection_id")
				fields := []models.ProjectCollectionField{}
				count, _ := app.Bun.NewSelect().Model(&fields).
					Where("collection_id = ?", collectionID).
					Order("created_at ASC").
					ScanAndCount(context.Background())
				paginatedResponse(ctx, fields, count)
			},
		},
		Show: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				field := &models.ProjectCollectionField{}
				err := app.Bun.NewSelect().Model(field).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Scan(context.Background())
				if err != nil || field.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Field not found")
					return
				}
				ctx.JSON(http.StatusOK, field)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				var form FieldForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Type == "" {
					form.Type = "text"
				}
				err := validation.ValidateStruct(&form,
					validation.Field(&form.Name, validation.Required),
				)
				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, app.FormatErrors(err))
					return
				}
				field := &models.ProjectCollectionField{
					ID:           xid.New().String(),
					CollectionID: ctx.Param("collection_id"),
					Name:         form.Name,
					Type:         form.Type,
					IsRequired:   form.IsRequired,
					IsIndexed:    form.IsIndexed,
					IsUnique:     form.IsUnique,
					IsSortable:   form.IsSortable,
					IsFilterable: form.IsFilterable,
					Options:      form.Options,
				}
				if field.Options == nil {
					field.Options = map[string]any{}
				}
				_, err = app.Bun.NewInsert().Model(field).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to create field")
					return
				}
				ctx.JSON(http.StatusOK, field)
			},
		},
		Update: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				field := &models.ProjectCollectionField{}
				err := app.Bun.NewSelect().Model(field).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Scan(context.Background())
				if err != nil || field.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Field not found")
					return
				}
				var form FieldForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Name != "" {
					field.Name = form.Name
				}
				if form.Type != "" {
					field.Type = form.Type
				}
				field.IsRequired = form.IsRequired
				field.IsIndexed = form.IsIndexed
				field.IsUnique = form.IsUnique
				field.IsSortable = form.IsSortable
				field.IsFilterable = form.IsFilterable
				if form.Options != nil {
					field.Options = form.Options
				}
				_, err = app.Bun.NewUpdate().Model(field).WherePK().Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to update field")
					return
				}
				ctx.JSON(http.StatusOK, field)
			},
		},
		Delete: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				_, err := app.Bun.NewDelete().Model((*models.ProjectCollectionField)(nil)).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to delete field")
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		},
	}
}
