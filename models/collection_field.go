package models

import "time"

type CollectionField struct {
	ID           string             `bun:"id,pk,notnull"`
	CollectionID int64              `bun:"collection_id,notnull"`
	Collection   *ProjectCollection `bun:"rel:belongs-to,join:collection_id=id"`

	Name         string         `bun:"name"`
	Type         string         `bun:"type,default:'base'"`
	IsRequired   bool           `bun:"is_required,default:0"`
	IsIndexed    bool           `bun:"is_indexed,default:0"`
	IsUnique     bool           `bun:"is_unique,default:0"`
	IsSortable   bool           `bun:"is_sortable,default:0"`
	IsFilterable bool           `bun:"is_filterable,default:0"`
	Options      map[string]any `bun:"options,type:json"`
	CreatedAt    time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
