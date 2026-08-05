package core

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Repository struct {
	Bun *bun.DB
}

func NewRepository(bun *bun.DB) *Repository {
	return &Repository{
		Bun: bun,
	}
}

func (r *Repository) GetById(result any, table string, id string) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1", table)
	r.Bun.NewRaw(sql, id).Scan(context.Background(), result)
}

func (r *Repository) GetByColumn(result any, table string, column string, value string) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1", table, column)
	r.Bun.NewRaw(sql, value).Scan(context.Background(), result)
}
