package handlers

import (
	"context"
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/xid"
)

type UserForm struct {
	ProjectID   *string        `json:"project_id" form:"project_id"`
	FirstName   string         `json:"first_name" form:"first_name"`
	LastName    string         `json:"last_name" form:"last_name"`
	Email       string         `json:"email" form:"email"`
	Password    string         `json:"password" form:"password"`
	IsActive    bool           `json:"is_active" form:"is_active"`
	Permissions map[string]any `json:"permissions" form:"permissions"`
}

func NewUsersHandler() *core.ResourceHandler {
	return &core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				users := []models.User{}
				query := app.Bun.NewSelect().Model(&users).Order("created_at DESC")
				if projectID := ctx.Query("project_id"); projectID != "" {
					query = query.Where("project_id = ?", projectID)
				}
				count, _ := query.ScanAndCount(context.Background())
				paginatedResponse(ctx, users, count)
			},
		},
		Show: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				user := &models.User{}
				err := app.Bun.NewSelect().Model(user).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || user.ID == "" {
					jsonError(ctx, http.StatusNotFound, "User not found")
					return
				}
				ctx.JSON(http.StatusOK, user)
			},
		},
		Create: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				var form UserForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				err := validation.ValidateStruct(&form,
					validation.Field(&form.FirstName, validation.Required),
					validation.Field(&form.Email, validation.Required, is.Email),
					validation.Field(&form.Password, validation.Required),
				)
				if err != nil {
					ctx.JSON(http.StatusUnprocessableEntity, app.FormatErrors(err))
					return
				}
				user := &models.User{
					ID:          xid.New().String(),
					ProjectID:   form.ProjectID,
					FirstName:   form.FirstName,
					LastName:    form.LastName,
					Email:       form.Email,
					Password:    utils.HashPassword(form.Password),
					IsActive:    form.IsActive,
					Permissions: form.Permissions,
				}
				if user.Permissions == nil {
					user.Permissions = map[string]any{}
				}
				_, err = app.Bun.NewInsert().Model(user).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to create user")
					return
				}
				user.Password = ""
				ctx.JSON(http.StatusOK, user)
			},
		},
		Update: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				user := &models.User{}
				err := app.Bun.NewSelect().Model(user).Where("id = ?", ctx.Param("id")).Scan(context.Background())
				if err != nil || user.ID == "" {
					jsonError(ctx, http.StatusNotFound, "User not found")
					return
				}
				var form UserForm
				if err := app.BindSafely(ctx, &form); err != nil {
					jsonError(ctx, http.StatusBadRequest, err.Error())
					return
				}
				if form.FirstName != "" {
					user.FirstName = form.FirstName
				}
				if form.LastName != "" {
					user.LastName = form.LastName
				}
				if form.Email != "" {
					user.Email = form.Email
				}
				if form.Password != "" {
					user.Password = utils.HashPassword(form.Password)
				}
				user.IsActive = form.IsActive
				if form.ProjectID != nil {
					user.ProjectID = form.ProjectID
				}
				if form.Permissions != nil {
					user.Permissions = form.Permissions
				}
				_, err = app.Bun.NewUpdate().Model(user).WherePK().Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to update user")
					return
				}
				user.Password = ""
				ctx.JSON(http.StatusOK, user)
			},
		},
		Delete: &core.ResourceAction{
			Handler: func(ctx *gin.Context, app *core.App) {
				if _, ok := getAuthUser(ctx, app); !ok {
					return
				}
				_, err := app.Bun.NewDelete().Model((*models.User)(nil)).Where("id = ?", ctx.Param("id")).Exec(context.Background())
				if err != nil {
					jsonError(ctx, http.StatusUnprocessableEntity, "Unable to delete user")
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"success": true})
			},
		},
	}
}
