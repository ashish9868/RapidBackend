package respository

import "github.com/uptrace/bun"

type AuthRepository struct {
	DB *bun.DB
}

func NewAuthRepository(db *bun.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func GetSuperAdminByEmail(email string) {

}
