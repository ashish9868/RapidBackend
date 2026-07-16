package models

import "time"

type Project struct {
	ID          string         `bun:"id,pk,notnull"`
	Title       string         `bun:"title,notnull"`
	Descriptiom string         `bun:"description"`
	Slug        string         `bun:"slug,unique,notnull"`
	Settings    map[string]any `bun:"type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
