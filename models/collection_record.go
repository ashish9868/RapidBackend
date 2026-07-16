package models

import "time"

type CollectionRecord struct {
	ID           string         `bun:"id,pk,notnull"`
	CollectionID int64          `bun:"collection_id,notnull"`
	Collection   *Collection    `bun:"rel:belongs-to,join:collection_id=id"`
	Data         map[string]any `bun:"data,type:json"`
	Version      int            `bun:"version,default:1"`
	Name         string         `bun:"name"`
	Type         string         `bun:"type,default:'base'"`
	IsRequired   bool           `bun:"is_required,default:0"`
	IsIndexed    bool           `bun:"is_indexed,default:0"`
	IsUnique     bool           `bun:"is_unique,default:0"`
	IsSortable   bool           `bun:"is_sortable,default:0"`
	IsFilterable bool           `bun:"is_filterable,default:0"`
	Options      map[string]any `bun:"options,type:json"`
	CreatedByID  int64          `bun:"created_by_id,notnull"`
	CreatedBy    *ProjectUser   `bun:"rel:belongs-to,join:created_by_id=id"`
	UpdatedByID  int64          `bun:"updated_by_id,notnull"`
	UpdatedBy    *ProjectUser   `bun:"rel:belongs-to,join:updated_by_id=id"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
