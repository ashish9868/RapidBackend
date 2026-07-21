package services

import (
	"context"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/rs/xid"
)

type AuthService struct {
	App *core.App
}

func NewAuthService(app *core.App) *AuthService {
	return &AuthService{App: app}
}

func (a *AuthService) LoginByEmail(email string, password string) *models.AccessKeyToken {
	token := a.LoginByEmailSuperadmin(email, password)
	if token != nil {
		return token
	}
	user := models.ProjectUser{}
	err := a.App.Bun.NewSelect().Model(&user).Where("email = ? AND email_verified_at IS NOT NULL AND is_active = ?", email, true).Scan(context.Background())
	if err != nil {
		println("Error is selecting user: ", err.Error())
		return nil
	}
	if len(user.ID) < 1 {
		return nil
	}
	if !a.App.BaseUtil.CheckPassword(user.Password, password) {
		return nil
	}

	expiry := time.Now().Add(1 * time.Hour)
	_, access_token, _ := a.App.BaseUtil.GenerateRandomHash()
	_, refresh_token, _ := a.App.BaseUtil.GenerateRandomHash()
	token = &models.AccessKeyToken{
		ID:           xid.New().String(),
		UserID:       &user.ID,
		ExpiresAt:    &expiry,
		Token:        access_token,
		CreatedAt:    time.Now(),
		RefreshToken: &refresh_token,
	}
	_, err = a.App.Bun.NewInsert().Model(token).Exec(context.Background())
	if err != nil {
		println("Error creating session: ", err.Error())
		return nil
	}

	return token
}

func (a *AuthService) LoginByEmailSuperadmin(email string, password string) *models.AccessKeyToken {
	superadmin := models.SuperAdmin{}
	err := a.App.Bun.NewSelect().Model(&superadmin).Where("email = ? AND email_verified_at IS NOT NULL AND is_active = ?", email, true).Scan(context.Background())
	if err != nil {
		println("Error is selecting user: ", err.Error())
		return nil
	}
	if len(superadmin.ID) < 1 {
		return nil
	}
	if !a.App.BaseUtil.CheckPassword(superadmin.Password, password) {
		return nil
	}
	expiry := time.Now().Add(1 * time.Hour)
	_, access_token, _ := a.App.BaseUtil.GenerateRandomHash()
	_, refresh_token, _ := a.App.BaseUtil.GenerateRandomHash()
	token := &models.AccessKeyToken{
		ID:           xid.New().String(),
		SuperAdminID: &superadmin.ID,
		ExpiresAt:    &expiry,
		Token:        access_token,
		CreatedAt:    time.Now(),
		RefreshToken: &refresh_token,
	}
	_, err = a.App.Bun.NewInsert().Model(token).Exec(context.Background())
	if err != nil {
		println("Error creating session: ", err.Error())
		return nil
	}

	return token
}
