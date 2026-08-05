package handlers

import (
	"context"
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type RecordForm struct {
	Data map[string]any `json:"data" form:"data"`
}

func NewRecordsHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				collectionID := ctx.Param("collection_id")
				records := []models.ProjectCollectionRecord{}
				count, _ := app.Bun.NewSelect().Model(&records).
					Where("collection_id = ?", collectionID).
					Order("created_at DESC").
					ScanAndCount(context.Background())
				paginatedResponse(ctx, records, count)
			},
		},
		Show: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				record := &models.ProjectCollectionRecord{}
				err := app.Bun.NewSelect().Model(record).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Scan(context.Background())
				if err != nil || record.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Record not found")
					return
				}
				ctx.JSON(http.StatusOK, record)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				authVal, ok := getAuthUser(ctx, app)
				if !ok {
					return
				}
				user := authVal.(*dto.AuthUser)
				var form RecordForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Data == nil {
					form.Data = map[string]any{}
				}
				record := &models.ProjectCollectionRecord{
					ID:           xid.New().String(),
					CollectionID: ctx.Param("collection_id"),
					Data:         form.Data,
					Version:      1,
					CreatedByID:  user.ID,
					UpdatedByID:  user.ID,
				}
				_, err := app.Bun.NewInsert().Model(record).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to create record")
					return
				}
				ctx.JSON(http.StatusOK, record)
			},
		},
		Update: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				authVal, ok := getAuthUser(ctx, app)
				if !ok {
					return
				}
				user := authVal.(*dto.AuthUser)
				record := &models.ProjectCollectionRecord{}
				err := app.Bun.NewSelect().Model(record).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Scan(context.Background())
				if err != nil || record.ID == "" {
					jsonError(ctx, http.StatusNotFound, "Record not found")
					return
				}
				var form RecordForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.Data != nil {
					record.Data = form.Data
				}
				record.Version++
				record.UpdatedByID = user.ID
				_, err = app.Bun.NewUpdate().Model(record).WherePK().Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to update record")
					return
				}
				ctx.JSON(http.StatusOK, record)
			},
		},
		Delete: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				_, err := app.Bun.NewDelete().Model((*models.ProjectCollectionRecord)(nil)).
					Where("id = ? AND collection_id = ?", ctx.Param("id"), ctx.Param("collection_id")).
					Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to delete record")
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		},
	}
}
