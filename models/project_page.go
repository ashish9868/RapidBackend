package models

import "time"

type ProjectPage struct {
	ID          string         `bun:"id,pk,notnull"`
	ProjectID   int64          `bun:"project_id,unique:unq_project_page,notnull"`
	Project     *Project       `bun:"rel:belongs-to,join:project_id=id"`
	Title       string         `bun:"title,notnull,unique:unq_project_page"`
	Descriptiom string         `bun:"description"`
	Slug        string         `bun:"slug,unique,notnull"`
	Settings    map[string]any `bun:"type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
