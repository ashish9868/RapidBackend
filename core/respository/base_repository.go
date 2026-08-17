package respository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type BaseRepository struct {
	DB *bun.DB
}

func NewBaseRepository(db *bun.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

func (b *BaseRepository) Escape(dbField string) string {
	parts := strings.Split(dbField, ".")
	for i, part := range parts {
		// Escape existing double quotes by doubling them
		clean := strings.ReplaceAll(part, `"`, `""`)
		parts[i] = `"` + clean + `"`
	}
	return strings.Join(parts, ".")
}
func (b *BaseRepository) GetByColumn(model any, column string, value string) error {
	return b.DB.NewSelect().Model(model).Where(fmt.Sprintf("%s = ?", b.Escape(column))).Scan(context.Background())
}

func (b *BaseRepository) GetById(model any, value string) error {
	return b.DB.NewSelect().Model(model).Where("id = ?", value).Scan(context.Background())
}

func (b *BaseRepository) DeleteById(model any, value string) (sql.Result, error) {
	return b.DB.NewDelete().Model(model).Where("id = ?", value).Exec(context.Background())
}
