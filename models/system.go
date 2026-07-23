package models

import "time"

type Role string

const (
	RoleSuperAdmin   Role = "superadmin"
	RoleProjectOwner Role = "project_owner"
	RoleProjectUser  Role = "project_user"
	RoleGuest        Role = "guest"
)

type AccessKeyToken struct {
	ID           string     `bun:"id,pk,notnull"`
	UserID       string     `bun:"project_user_id,notnull"`
	User         *User      `bun:"rel:belongs-to,join:project_user_id=id"`
	Token        string     `bun:"access_token"`
	RefreshToken string     `bun:"refresh_token"`
	CreatedAt    time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	ExpiresAt    *time.Time `bun:"expires_at,notnull"`
}

type Project struct {
	ID          string         `bun:"id,pk,notnull"`
	Title       string         `bun:"title,notnull"`
	Descriptiom string         `bun:"description"`
	Slug        string         `bun:"slug,unique,notnull"`
	Settings    map[string]any `bun:"type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type User struct {
	ID string `bun:"id,pk,notnull"`

	ProjectID *string  `bun:"project_id,unique:unq_project_email"`
	Project   *Project `bun:"rel:belongs-to,join:project_id=id"`

	FirstName       string     `bun:"first_name,notnull"`
	LastName        string     `bun:"last_name"`
	Email           string     `bun:"email,unique:unq_project_email,notnull"`
	Password        string     `bun:"password,notnull" json:"-"`
	EmailVerifiedAt *time.Time `bun:"email_verified_at,notnull"`
	IsActive        bool       `bun:"is_active,default:'0'"`
	Role            Role       `bun:"role,notnull,unique:unq_project_email,default:'owner'"`

	Permissions map[string]any `bun:"permissions_json,type:json"`
	CreatedAt   time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type ProjectCollection struct {
	ID string `bun:"id,pk,notnull"`

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
	ID           string             `bun:"id,pk,notnull"`
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
	ID           string             `bun:"id,pk,notnull"`
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
	ID          string         `bun:"id,pk,notnull"`
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
	ID               string    `bun:"id,pk,notnull"`
	Name             string    `bun:"name,unique"`
	Description      string    `bun:"description"`
	isSystemTemplate string    `bun:"is_system_template,default:1"`
	HtmlContent      string    `bun:"html_content"`
	TextContent      string    `bun:"text_content"`
	CreatedAt        time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type SystemSetting struct {
	ID          string    `bun:"id,pk,notnull"`
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
