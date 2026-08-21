package repository

import (
	"time"

	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
)

type AccessTokenRepository struct {
	BaseRepository *BaseRepository
}

func NewAccessTokenRepository(baseRepository *BaseRepository) *AccessTokenRepository {
	return &AccessTokenRepository{BaseRepository: baseRepository}
}

func (a *AccessTokenRepository) CreateNewAccessToken(user *models.User) *models.AccessKeyToken {
	token := &models.AccessKeyToken{
		CollectionID: user.ID,
		Collection:   utils.IfElse(user.Collection == COLLECTION_SUPERADMINS, COLLECTION_SUPERADMINS, COLLECTION_USERS),
		Token:        utils.GenerateRandomHash(),
		CreatedAt:    time.Now(),
		User:         user,
	}
	_, err := a.BaseRepository.Insert(COLLECTION_ACCESS_KEY_TOKENS, token)
	if err == nil {
		return token
	}

	utils.Log(err.Error())
	return nil
}

func (a *AccessTokenRepository) GetUserFromToken(accessToken string) *models.User {
	if len(accessToken) < 1 {
		return nil
	}
	token := &models.AccessKeyToken{}
	err := a.BaseRepository.GetByColumn(COLLECTION_ACCESS_KEY_TOKENS, "access_token", accessToken, token)
	if err != nil || !(token.ID > 0) || !(token.CollectionID > 0) {
		return nil
	}
	user := &models.User{}
	a.BaseRepository.GetById(token.Collection, user, utils.ToString(token.CollectionID))
	if user.ID > 0 {
		return user
	}
	return nil
}
