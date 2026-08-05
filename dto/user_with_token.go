package dto

import "time"

type AuthUser struct {
	ID              string
	ProjectID       *string
	Token           *string
	FirstName       string
	LastName        string
	IsSuperadmin    bool
	Email           string
	EmailVerifiedAt *time.Time
	Permissions     map[string]any
}
