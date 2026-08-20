package repository

import (
	"time"

	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/rs/xid"
)

type AccessTokenRepository struct {
	BaseRepository *BaseRepository
}

func NewAccessTokenRepository(baseRepository *BaseRepository) *AccessTokenRepository {
	return &AccessTokenRepository{BaseRepository: baseRepository}
}

func (a *AccessTokenRepository) CreateNewAccessToken(user *models.User) *models.AccessKeyToken {
	token := &models.AccessKeyToken{
		ID:           xid.New().String(),
		CollectionID: user.ID,
		Collection:   utils.IFElse(user.Collection == COLLECTION_SUPERADMINS, COLLECTION_SUPERADMINS, COLLECTION_USERS),
		Token:        utils.GenerateRandomHash(),
		CreatedAt:    time.Now(),
	}
	_, err := a.BaseRepository.Insert(token.Collection, token.ToMap())
	if err != nil {
		return token
	}
	return nil
}
