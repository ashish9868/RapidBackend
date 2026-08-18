package respository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ashish9868/rapidbackend/utils"
	"github.com/uptrace/bun"
)

type BaseRepository struct {
	DB *bun.DB
}

func NewBaseRepository(db *bun.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

func (b *BaseRepository) Escape(s string) string {
	parts := strings.Split(s, ".")
	invalidIdent := regexp.MustCompile(`[^a-zA-Z0-9_.]`)

	finalParts := []string{}

	for _, part := range parts {
		finalParts = append(finalParts, invalidIdent.ReplaceAllString(part, ""))
	}
	return strings.Join(finalParts, ".")
}
func (b *BaseRepository) BuildWhere(where map[string]any, operation string) (string, []any) {
	keys := []string{}
	values := []any{}
	for key, value := range where {
		keys = append(keys, fmt.Sprintf("%s %s = ?", operation, b.Escape(key)))
		values = append(values, value)
	}
	if len(keys) > 0 {
		keys = append([]string{`"1" = "1"`}, keys...)
	}
	utils.LogF("WHERE was %s", strings.Join(keys, " "))
	utils.Log(utils.ToJSON(values))
	return strings.Join(keys, " "), values
}

func (b *BaseRepository) GetByColumn(model any, column string, value string) error {
	return b.DB.NewSelect().Model(model).Where(fmt.Sprintf("%s = ?", b.Escape(column))).Scan(context.Background())
}

func (b *BaseRepository) GetById(model any, value string) error {
	return b.DB.NewSelect().Model(model).Where("id = ?", value).Scan(context.Background())
}

func (b *BaseRepository) SelectWhere(model any, where map[string]any) error {
	whereString, whereValue := b.BuildWhere(where, "AND")
	if len(whereString) > 0 && len(whereValue) > 0 {
		return b.DB.NewSelect().Model(model).Where(whereString, whereValue...).Scan(context.Background())
	}
	return errors.New("SelectWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) Exists(table, column string, value any, ignore_id string) bool {

	var result string
	var id string
	err := b.DB.NewSelect().
		Table(b.Escape(table)).
		ColumnExpr(`?`, b.Escape(column)).
		ColumnExpr(`?`, "id").
		Limit(1).
		Where(fmt.Sprintf("%s = ?", b.Escape(column)), value).
		Scan(context.Background(), &result, &id)
	if err == nil {
		if len(ignore_id) > 0 {
			return strings.EqualFold(ignore_id, id)
		}
		return true
	}
	return false
}

func (b *BaseRepository) DeleteById(model any, value string) (sql.Result, error) {
	return b.DB.NewDelete().Model(model).Where("id = ?", value).Exec(context.Background())
}

func (b *BaseRepository) Insert(model any) (sql.Result, error) {
	return b.DB.NewInsert().Model(model).Exec(context.Background())
}

func (b *BaseRepository) InsertOrUpdate(model any, where map[string]any) (sql.Result, error) {
	err := b.SelectWhere(model, where)
	if err != nil {
		return b.Insert(model)
	}
	return b.UpdateWhere(model, where)
}

func (b *BaseRepository) UpdateById(model any, id string) (sql.Result, error) {
	return b.DB.NewUpdate().Model(model).Where("id = ?", id).Exec(context.Background())
}

func (b *BaseRepository) UpdateWhere(model any, where map[string]any) (sql.Result, error) {
	whereString, whereValue := b.BuildWhere(where, "AND")
	if len(whereString) > 0 && len(whereValue) > 0 {
		return b.DB.NewUpdate().Model(model).Where(whereString, whereValue...).Exec(context.Background())
	}
	return nil, errors.New("UpdateWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) DeleteWhere(model any, where map[string]any) (sql.Result, error) {
	whereKey, whereValue := b.BuildWhere(where, "AND")
	if len(whereKey) > 0 && len(whereValue) > 0 {
		return b.DB.NewDelete().Model(model).Where(whereKey, whereValue...).Exec(context.Background())
	}
	return nil, errors.New("DeleteWhere: Operation not allowed with empty where")
}

func (app *BaseRepository) WithTransaction(ctx context.Context, db *bun.DB, fn func(tx bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (b *BaseRepository) GetByColumnTx(tx bun.Tx, model any, column string, value string) error {
	return tx.NewSelect().Model(model).Where(fmt.Sprintf("%s = ?", b.Escape(column))).Scan(context.Background())
}

func (b *BaseRepository) GetByIdTx(tx bun.Tx, model any, value string) error {
	return tx.NewSelect().Model(model).Where("id = ?", value).Scan(context.Background())
}

func (b *BaseRepository) SelectWhereTx(tx bun.Tx, model any, where map[string]any) error {
	whereString, whereValue := b.BuildWhere(where, "AND")
	return tx.NewSelect().Model(model).Where(whereString, whereValue...).Scan(context.Background())
}

func (b *BaseRepository) DeleteByIdTx(tx bun.Tx, model any, value string) (sql.Result, error) {
	return tx.NewDelete().Model(model).Where("id = ?", value).Exec(context.Background())
}

func (b *BaseRepository) InsertTx(tx bun.Tx, model any) (sql.Result, error) {
	return tx.NewInsert().Model(model).Exec(context.Background())
}

func (b *BaseRepository) UpdateByIdTx(tx bun.Tx, model any, id string) (sql.Result, error) {
	return tx.NewUpdate().Model(model).Where("id = ?", id).Exec(context.Background())
}

func (b *BaseRepository) InsertOrUpdateTx(tx bun.Tx, model any, where map[string]any) (sql.Result, error) {
	err := b.SelectWhereTx(tx, model, where)
	if err == nil {
		return b.InsertTx(tx, model)
	}
	return b.UpdateWhere(model, where)
}

func (b *BaseRepository) UpdateWhereTx(tx bun.Tx, model any, where map[string]any) (sql.Result, error) {
	whereKey, whereValue := b.BuildWhere(where, "AND")
	if len(whereKey) > 0 && len(whereValue) > 0 {
		return tx.NewUpdate().Model(model).Where(whereKey, whereValue...).Exec(context.Background())
	}
	return nil, errors.New("UpdateWhereTx: Operation not allowed with empty where")
}

func (b *BaseRepository) DeleteWhereTx(tx bun.Tx, model any, where map[string]any) (sql.Result, error) {
	whereKey, whereValue := b.BuildWhere(where, "AND")
	if len(whereKey) > 0 && len(whereValue) > 0 {
		return tx.NewDelete().Model(model).Where(whereKey, whereValue...).Exec(context.Background())
	}
	return nil, errors.New("DeleteWhereTx: Operation not allowed with empty where")
}
