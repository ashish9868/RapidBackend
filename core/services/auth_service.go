package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type AuthService struct {
	Bun *bun.DB
}

func NewAuthService(bun *bun.DB) *AuthService {
	return &AuthService{Bun: bun}
}

func (a *AuthService) LoginByEmail(email string, password string, collection string) *dto.AuthUser {

	var ID string
	var Password string

	sql := fmt.Sprintf(`SELECT id, password, email FROM %s WHERE email = ?`, collection)
	err := a.Bun.NewRaw(sql, email).Scan(context.Background(), &ID, &Password)

	if len(ID) < 1 {
		return nil
	}
	if !utils.CheckPassword(Password, password) {
		return nil
	}
	_, access_token, _ := utils.GenerateRandomHash()

	user := a.GetAnyUserById(collection, ID)

	hashedToken := utils.HashPassword(access_token)
	token := &models.AccessKeyToken{
		ID:           xid.New().String(),
		CollectionID: user.ID,
		Collection:   collection,
		Token:        hashedToken,
		CreatedAt:    time.Now(),
	}
	_, err = a.Bun.NewInsert().Model(token).Exec(context.Background())
	if err != nil {
		println("Error creating session: ", err.Error())
		return nil
	}

	user.Token = &hashedToken
	return user
}

func (a *AuthService) GetUserByToken(token string) *dto.AuthUser {
	if len(token) < 1 {
		return nil
	}
	aToken := &models.AccessKeyToken{}
	err := a.Bun.NewSelect().Model(aToken).Where("access_token = ?", token).Scan(context.Background())
	if err != nil {
		fmt.Println("Error while testing token: ", err.Error())
		return nil
	}
	if aToken.Collection == "superadmins" {
		return a.GetSuperAdminById(aToken.CollectionID)
	}
	return a.GetUserById(aToken.CollectionID)
}

func (a *AuthService) GetAnyUserById(collection string, id string) *dto.AuthUser {
	if collection == "superadmins" {
		return a.GetSuperAdminById(id)
	}
	return a.GetUserById(id)
}

func (a *AuthService) GetUserById(id string) *dto.AuthUser {
	user := &models.User{}
	if len(strings.TrimSpace(id)) > 0 {
		a.Bun.NewSelect().Model(user).Where("id = ?", id).Scan(context.Background())
	}
	if user.ID == id {
		return &dto.AuthUser{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			IsSuperadmin:    false,
			EmailVerifiedAt: nil,
		}

	}
	return nil
}

func (a *AuthService) GetSuperAdminById(id string) *dto.AuthUser {
	user := &models.Superadmin{}
	if len(strings.TrimSpace(id)) > 0 {
		a.Bun.NewSelect().Model(user).Where("id = ?", id).Scan(context.Background())
	}
	if user.ID == id {
		return &dto.AuthUser{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			IsSuperadmin:    true,
			EmailVerifiedAt: nil,
		}
	}
	return nil
}
