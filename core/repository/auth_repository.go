package repository

import (
	"github.com/ashish9868/rapidbackend/models"
)

type AuthRepository struct {
	BaseRepository *BaseRepository
}

func NewAuthRepository(baseRepository *BaseRepository) *AuthRepository {
	return &AuthRepository{BaseRepository: baseRepository}
}

func (a *AuthRepository) GetUserByToken(token string) *models.User {
	return nil
}
