package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type BaseUtil struct{}

func NewBaseUtil() *BaseUtil {
	return &BaseUtil{}
}

func (util *BaseUtil) SafeEnvGet(name string, defaultVal string) string {
	val := strings.TrimSpace(os.Getenv(name))
	if len(val) > 0 {
		return val
	}
	return defaultVal
}

func (uitl *BaseUtil) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return errors.New("password must contain an uppercase letter")
	case !hasLower:
		return errors.New("password must contain a lowercase letter")
	case !hasDigit:
		return errors.New("password must contain a number")
	case !hasSpecial:
		return errors.New("password must contain a special character")
	}

	return nil
}

func (BaseUtil) HashPassword(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}
func (BaseUtil) CheckPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func (BaseUtil) GetSha256Hash(str string) string {
	h := sha256.Sum256([]byte(str))
	return hex.EncodeToString(h[:])
}
