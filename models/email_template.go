package models

import "time"

type EmailTemplate struct {
	ID               string    `bun:"id,pk,notnull"`
	Name             string    `bun:"name,unique"`
	Description      string    `bun:"description"`
	isSystemTemplate string    `bun:"is_system_template,default:1"`
	HtmlContent      string    `bun:"html_content"`
	TextContent      string    `bun:"text_content"`
	CreatedAt        time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
