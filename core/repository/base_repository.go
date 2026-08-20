package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ashish9868/rapidbackend/utils"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

const (
	COLLECTION_SUPERADMINS                = "superadmins"
	COLLECTION_USERS                      = "users"
	COLLECTION_ACCESS_KEY_TOKENS          = "access_key_tokens"
	COLLECTION_PROJECTS                   = "projects"
	COLLECTION_PROJECT_COLLECTIONS        = "project_collections"
	COLLECTION_PROJECT_PAGES              = "project_pages"
	COLLECTION_PROJECT_COLLECTION_FIELDS  = "project_collection_fields"
	COLLECTION_PROJECT_COLLECTION_RECORDS = "project_collection_records"
	COLLECTION_SETTINGS                   = "settings"
)

type BaseRepository struct {
	DB *sqlx.DB
}

func NewBaseRepository(db *sqlx.DB) *BaseRepository {
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
func (b *BaseRepository) BuildWhere(where map[string]any) (placeholders []string, values []any) {
	placeholders = make([]string, 0, len(where))
	values = make([]any, 0, len(where))
	for key, value := range where {
		placeholders = append(placeholders, fmt.Sprintf("%s = ?", b.Escape(key)))
		values = append(values, value)
	}
	return placeholders, values
}

func (b *BaseRepository) BuildInsertPlaceHolders(data map[string]any) (placeholders []string, columns []string, values []any) {
	columns = make([]string, 0, len(data))
	values = make([]any, 0, len(data))
	placeholders = make([]string, 0, len(data))
	for column, value := range data {
		columns = append(columns, b.Escape(column))
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}
	return placeholders, columns, values
}

func (b *BaseRepository) BuildUpdatePlaceHolders(data map[string]any) (placeholders []string, values []any) {
	values = make([]any, 0, len(data))
	placeholders = make([]string, 0, len(data))
	for key, value := range data {
		placeholders = append(placeholders, fmt.Sprintf("%s = ?", b.Escape(key)))
		values = append(values, value)
	}
	return placeholders, values

}

func (b *BaseRepository) ModelToMap(model any) map[string]any {
	return utils.ModelToMap(model, "db")
}

func (b *BaseRepository) GetByColumn(table string, column string, value string, model any) error {
	placeholders, values := b.BuildWhere(map[string]any{
		column: value,
	})
	sql := fmt.Sprintf(
		`SELECT * FROM %s WHERE %s LIMIT 1`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
	utils.Log("SQL: GetByColumn -> ", utils.ToJSON(map[string]any{
		"query":        sql,
		"placeholders": placeholders,
		"values":       values,
	}))

	return b.DB.GetContext(context.Background(), model, sql, values...)
}

func (b *BaseRepository) GetById(table string, model any, value string) error {
	return b.GetByColumn(table, "id", value, model)
}

func (b *BaseRepository) SelectWhere(table string, model any, where map[string]any) error {
	placeholders, values := b.BuildWhere(where)
	modelMap := b.ModelToMap(model)
	if len(placeholders) > 0 && len(values) > 0 {
		sql := fmt.Sprintf(
			`SELECT "%s" FROM %s WHERE %s LIMIT 1`,
			strings.Join(lo.Keys(modelMap), `","`),
			b.Escape(table),
			strings.Join(placeholders, " AND "),
		)
		utils.Log("SQL: SelectWhere -> ", utils.ToJSON(map[string]any{
			"query":        sql,
			"placeholders": placeholders,
			"values":       values,
		}))
		return b.DB.GetContext(context.Background(), model, sql, values...)

	}
	return errors.New("SelectWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) DeleteWhere(table string, where map[string]any, limit int) (sql.Result, error) {
	placeholders, values := b.BuildWhere(where)
	if len(placeholders) > 0 && len(values) > 0 {
		sql := fmt.Sprintf(
			`DELETE FROM %s WHERE id IN(
				SELECT id FROM %s WHERE %s LIMIT %d
			)`,
			b.Escape(table),
			b.Escape(table),
			strings.Join(placeholders, " AND "),
			limit,
		)
		utils.Log("SQL: DeleteWhere -> ", utils.ToJSON(map[string]any{
			"query":        sql,
			"placeholders": placeholders,
			"values":       values,
		}))

		return b.DB.ExecContext(context.Background(), sql, values...)
	}
	return nil, errors.New("SelectWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) Exists(table, column string, value any, ignore_id string) bool {
	query := fmt.Sprintf(`
		SELECT %s, id
		FROM %s
		WHERE %s = ?
		LIMIT 1
	`,
		b.Escape(column),
		b.Escape(table),
		b.Escape(column),
	)

	utils.Log("SQL: Exists -> ", utils.ToJSON(map[string]any{
		"query": query,
	}))

	err := b.DB.GetContext(
		context.Background(),
		&struct {
			Result string `db:"result"`
			ID     int64  `db:"id"`
		}{},
		query,
		value,
	)

	// Easier with a temporary struct:
	return err == nil
}

func (b *BaseRepository) DeleteById(table, id string) (sql.Result, error) {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = ?`,
		b.Escape(table),
	)

	utils.Log("SQL: DeleteById -> ", utils.ToJSON(map[string]any{
		"query": query,
	}))

	return b.DB.ExecContext(
		context.Background(),
		query,
		id,
	)
}

func (b *BaseRepository) Insert(table string, model any) (sql.Result, error) {
	modelMap := lo.OmitByKeys(b.ModelToMap(model), []string{"id"})
	placeholders, columns, values := b.BuildInsertPlaceHolders(modelMap)
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		b.Escape(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	utils.Log("SQL: Insert -> ", utils.ToJSON(map[string]any{
		"query":        query,
		"placeholders": placeholders,
		"values":       values,
		"columns":      columns,
	}))

	return b.DB.ExecContext(
		context.Background(),
		query,
		values...,
	)
}

func (b *BaseRepository) InsertOrUpdate(table string, model any, where map[string]any) error {
	// check exists
	selectModel := &struct {
		ID int64 `db:"id"`
	}{}

	return b.WithTransaction(context.Background(), b.DB, func(tx *sqlx.Tx) error {
		err := b.SelectWhereTx(tx, table, selectModel, where)
		if err != nil {
			_, err = b.InsertTx(tx, table, model)
			if err != nil {
				return err
			}
		}
		_, err = b.UpdateByIdTx(tx, table, selectModel.ID, model)
		if err != nil {
			return err
		}
		return nil
	})

}

func (b *BaseRepository) UpdateById(table string, id int64, model any) (sql.Result, error) {
	modelMap := b.ModelToMap(model)
	placeholders, values := b.BuildUpdatePlaceHolders(modelMap)
	query := fmt.Sprintf(
		`UPDATE %s SET %s  WHERE id = ?`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)

	utils.Log("SQL: UpdateById -> ", utils.ToJSON(map[string]any{
		"query":        query,
		"placeholders": placeholders,
		"values":       values,
	}))

	return b.DB.ExecContext(
		context.Background(),
		query,
		append([]any{id}, values...),
	)
}

func (b *BaseRepository) GetByColumnTx(tx *sqlx.Tx, table string, column string, value string, model any) error {
	placeholders, values := b.BuildWhere(map[string]any{
		column: value,
	})
	sql := fmt.Sprintf(
		`SELECT * FROM %s WHERE %s LIMIT 1`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
	utils.Log("SQL: GetByColumnTx -> ", utils.ToJSON(map[string]any{
		"query":        sql,
		"placeholders": placeholders,
		"values":       values,
	}))

	return tx.GetContext(context.Background(), model, sql, values...)
}

func (b *BaseRepository) GetByIdTx(tx *sqlx.Tx, table string, model any, value string) error {
	return b.GetByColumnTx(tx, table, "id", value, model)
}

func (b *BaseRepository) SelectWhereTx(tx *sqlx.Tx, table string, model any, where map[string]any) error {
	placeholders, values := b.BuildWhere(where)
	modelMap := b.ModelToMap(model)
	if len(placeholders) > 0 && len(values) > 0 {
		sql := fmt.Sprintf(
			`SELECT "%s" FROM %s WHERE %s LIMIT 1`,
			strings.Join(lo.Keys(modelMap), `","`),
			b.Escape(table),
			strings.Join(placeholders, " AND "),
		)
		utils.Log("SQL: SelectWhereTx -> ", utils.ToJSON(map[string]any{
			"query":        sql,
			"placeholders": placeholders,
			"values":       values,
		}))

		return tx.GetContext(context.Background(), model, sql, values...)

	}
	return errors.New("SelectWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) DeleteByIdTx(tx *sqlx.Tx, table, id string) (sql.Result, error) {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = ?`,
		b.Escape(table),
	)

	utils.Log("SQL: DeleteByIdTx -> ", utils.ToJSON(map[string]any{
		"query": query,
	}))

	return tx.ExecContext(
		context.Background(),
		query,
		id,
	)
}

func (b *BaseRepository) InsertTx(tx *sqlx.Tx, table string, model any) (sql.Result, error) {

	modelMap := lo.OmitByKeys(b.ModelToMap(model), []string{"id"})
	placeholders, columns, values := b.BuildInsertPlaceHolders(modelMap)

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		b.Escape(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	utils.Log("SQL: InsertTx -> ", utils.ToJSON(map[string]any{
		"query":        query,
		"placeholders": placeholders,
		"values":       values,
		"columns":      columns,
	}))

	return tx.ExecContext(
		context.Background(),
		query,
		values...,
	)
}

func (b *BaseRepository) InsertOrUpdateTx(tx *sqlx.Tx, table string, model any, where map[string]any) (sql.Result, error) {
	// check exists
	selectModel := &struct {
		ID int64 `db:"id"`
	}{}

	modelMap := b.ModelToMap(model)

	err := b.SelectWhereTx(tx, table, &selectModel, where)
	if err != nil {
		return b.InsertTx(tx, table, modelMap)
	}
	return b.UpdateByIdTx(tx, table, selectModel.ID, model)
}

func (b *BaseRepository) UpdateByIdTx(tx *sqlx.Tx, table string, id int64, model any) (sql.Result, error) {
	modelMap := b.ModelToMap(model)
	placeholders, values := b.BuildUpdatePlaceHolders(modelMap)
	query := fmt.Sprintf(
		`UPDATE %s SET %s  WHERE id = ?`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
	utils.Log("SQL: UpdateByIdTx -> ", utils.ToJSON(map[string]any{
		"query":        query,
		"placeholders": placeholders,
		"values":       values,
	}))
	return tx.ExecContext(
		context.Background(),
		query,
		append(values, id)...,
	)
}

func (app *BaseRepository) WithTransaction(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
