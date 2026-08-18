package utils

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
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
