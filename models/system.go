package models

import (
	"time"
)

type AccessKeyToken struct {
	ID           int64     `db:"id" json:"id"`
	CollectionID int64     `db:"collection_id" json:"-"`
	Collection   string    `db:"collection" json:"-"`
	Token        string    `db:"access_token" json:"access_token"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	User         *User     `json:"user"`
}

type Project struct {
	ID          int64          `bun:"id,pk,notnull"`
	Title       string         `bun:"title,notnull"`
	Descriptiom string         `bun:"description"`
	Slug        string         `bun:"slug,unique,notnull"`
	Settings    map[string]any `bun:"type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type User struct {
	ID              int64          `db:"id" json:"id"`
	ProjectID       *string        `db:"project_id" json:"project_id"`
	FirstName       *string        `db:"first_name" json:"first_name"`
	LastName        *string        `db:"last_name" json:"last_name"`
	Email           string         `db:"email" json:"email"`
	Password        string         `db:"password" json:"-"`
	EmailVerifiedAt *time.Time     `db:"email_verified_at" json:"-"`
	IsActive        bool           `db:"is_active" json:"is_active"`
	Collection      string         `json:"-"`
	Permissions     map[string]any `db:"permissions_json" json:"permissions"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

type ProjectCollection struct {
	ID int64 `bun:"id,pk,notnull"`

	ProjectID string   `bun:"project_id,notnull"`
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

type ProjectCollectionField struct {
	ID           int64              `bun:"id,pk,notnull"`
	CollectionID string             `bun:"collection_id,notnull"`
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

type ProjectCollectionRecord struct {
	ID           int64              `bun:"id,pk,notnull"`
	CollectionID string             `bun:"collection_id,notnull"`
	Collection   *ProjectCollection `bun:"rel:belongs-to,join:collection_id=id"`
	Data         map[string]any     `bun:"data,type:json"`
	Version      int                `bun:"version,default:1"`
	Name         string             `bun:"name"`
	Type         string             `bun:"type,default:'base'"`
	IsRequired   bool               `bun:"is_required,default:0"`
	IsIndexed    bool               `bun:"is_indexed,default:0"`
	IsUnique     bool               `bun:"is_unique,default:0"`
	IsSortable   bool               `bun:"is_sortable,default:0"`
	IsFilterable bool               `bun:"is_filterable,default:0"`
	Options      map[string]any     `bun:"options,type:json"`
	CreatedByID  string             `bun:"created_by_id,notnull"`
	CreatedBy    *User              `bun:"rel:belongs-to,join:created_by_id=id"`
	UpdatedByID  string             `bun:"updated_by_id,notnull"`
	UpdatedBy    *User              `bun:"rel:belongs-to,join:updated_by_id=id"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type ProjectPage struct {
	ID          int64          `bun:"id,pk,notnull"`
	ProjectID   string         `bun:"project_id,unique:unq_project_page,notnull"`
	Project     *Project       `bun:"rel:belongs-to,join:project_id=id"`
	Title       string         `bun:"title,notnull,unique:unq_project_page"`
	Descriptiom string         `bun:"description"`
	Slug        string         `bun:"slug,unique,notnull"`
	Settings    map[string]any `bun:"type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type EmailTemplate struct {
	ID               int64     `bun:"id,pk,notnull"`
	Name             string    `bun:"name,unique"`
	Description      string    `bun:"description"`
	isSystemTemplate string    `bun:"is_system_template,default:1"`
	HtmlContent      string    `bun:"html_content"`
	TextContent      string    `bun:"text_content"`
	CreatedAt        time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type SystemSetting struct {
	ID          int64     `bun:"id,pk,notnull"`
	ProjectID   *string   `bun:"project_id"`
	Project     *Project  `bun:"rel:belongs-to,join:project_id=id"`
	Value       string    `bun:"value"`
	CreatedByID string    `bun:"created_by_id,notnull"`
	CreatedBy   *User     `bun:"rel:belongs-to,join:created_by_id=id"`
	UpdatedByID string    `bun:"updated_by_id,notnull"`
	UpdatedBy   *User     `bun:"rel:belongs-to,join:updated_by_id=id"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
