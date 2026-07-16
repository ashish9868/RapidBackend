package models

import "time"

type SuperAdmin struct {
	ID              string     `bun:"id,pk,notnull"`
	FirstName       string     `bun:"first_name,notnull"`
	LastName        string     `bun:"last_name"`
	Email           string     `bun:"email,unique,notnull"`
	Password        string     `bun:"password,notnull"`
	EmailVerifiedAt *time.Time `bun:"email_verified_at,notnull"`
	IsActive        bool       `bun:"is_active,default:0"`
	CreatedAt       time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt       time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
