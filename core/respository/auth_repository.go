package respository

import (
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
)

type AuthRepository struct {
	BaseRepository *BaseRepository
}

func NewAuthRepository(baseRepository *BaseRepository) *AuthRepository {
	return &AuthRepository{BaseRepository: baseRepository}
}

func (a *AuthRepository) GetSuperAdminBy(column string, value string) *models.Superadmin {
	user := &models.Superadmin{}
	err := a.BaseRepository.GetByColumn(user, column, value)
	if err != nil {
		utils.Log("Unable to find superadmin", err.Error())
		return nil
	}
	return user
}

func (a *AuthRepository) GetUserBy(column string, value string) *models.User {
	user := &models.User{}
	err := a.BaseRepository.GetByColumn(user, column, value)
	if err != nil {
		utils.Log("Unable to find user", err.Error())
		return nil
	}
	return user
}
