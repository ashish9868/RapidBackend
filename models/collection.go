package models

import "time"

type Collection struct {
	ID string `bun:"id,pk,notnull"`

	ProjectID int64    `bun:"project_id,notnull"`
	Project   *Project `bun:"rel:belongs-to,join:project_id=id"`

	Name      string         `bun:"name"`
	Type      string         `bun:"type,default:'base'"`
	SortOrder int            `bun:"sort_order,default:0"`
	Required  bool           `bun:"required,default:0"`
	Rules     map[string]any `bun:"rules,type:json"`
	Options   map[string]any `bun:"options,type:json"`
	CreatedAt time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
