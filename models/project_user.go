package models

import "time"

type ProjectUser struct {
	ID string `bun:"id,pk,notnull"`

	ProjectID int64    `bun:"project_id,unique:unq_project_email,notnull"`
	Project   *Project `bun:"rel:belongs-to,join:project_id=id"`

	FirstName       string     `bun:"first_name,notnull"`
	LastName        string     `bun:"last_name"`
	Email           string     `bun:"email,unique:unq_project_email,notnull"`
	Password        string     `bun:"password,notnull"`
	EmailVerifiedAt *time.Time `bun:"email_verified_at,notnull"`
	IsActive        bool       `bun:"is_active,default:'0'"`
	Role            string     `bun:"role,notnull,default:'owner'"`

	Permissions map[string]any `bun:"permissions_json,type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
