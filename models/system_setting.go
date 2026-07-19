package models

import "time"

type SystemSetting struct {
	ID          string       `bun:"id,pk,notnull"`
	ProjectID   *int64       `bun:"project_id"`
	Project     *Project     `bun:"rel:belongs-to,join:project_id=id"`
	Value       string       `bun:"value"`
	CreatedByID int64        `bun:"created_by_id,notnull"`
	CreatedBy   *ProjectUser `bun:"rel:belongs-to,join:created_by_id=id"`
	UpdatedByID int64        `bun:"updated_by_id,notnull"`
	UpdatedBy   *ProjectUser `bun:"rel:belongs-to,join:updated_by_id=id"`
	CreatedAt   time.Time    `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time    `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
