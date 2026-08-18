package respository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
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

func (b *BaseRepository) GetByColumn(table string, column string, value string, model any) error {
	placeholders, values := b.BuildWhere(map[string]any{
		column: value,
	})
	sql := fmt.Sprintf(
		`SELECT * FROM %s WHERE %s LIMIT 1`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
	return b.DB.GetContext(context.Background(), model, sql, values...)
}

func (b *BaseRepository) GetById(table string, model any, value string) error {
	return b.GetByColumn(table, "id", value, model)
}

func (b *BaseRepository) SelectWhere(table string, model any, where map[string]any) error {
	placeholders, values := b.BuildWhere(where)

	if len(placeholders) > 0 && len(values) > 0 {
		sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", b.Escape(table), strings.Join(placeholders, " AND "))
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

	err := b.DB.GetContext(
		context.Background(),
		&struct {
			Result string `db:"result"`
			ID     string `db:"id"`
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

	return b.DB.ExecContext(
		context.Background(),
		query,
		id,
	)
}

func (b *BaseRepository) Insert(table string, data map[string]any) (sql.Result, error) {

	placeholders, columns, values := b.BuildInsertPlaceHolders(data)

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		b.Escape(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return b.DB.ExecContext(
		context.Background(),
		query,
		values...,
	)
}

func (b *BaseRepository) InsertOrUpdate(table string, data map[string]any, where map[string]any) error {
	// check exists
	model := struct {
		Result string `db:"result"`
		ID     string `db:"id"`
	}{}

	return b.WithTransaction(context.Background(), b.DB, func(tx *sqlx.Tx) error {
		err := b.SelectWhere(table, model, where)
		if err != nil {
			_, err = b.InsertTx(tx, table, data)
			if err != nil {
				return err
			}
		}
		_, err = b.UpdateByIdTx(tx, table, model.ID, data)
		if err != nil {
			return err
		}
		return nil
	})

}

func (b *BaseRepository) UpdateById(table string, id string, data map[string]any) (sql.Result, error) {
	placeholders, values := b.BuildUpdatePlaceHolders(data)
	query := fmt.Sprintf(
		`UPDATE %s SET %s  WHERE id = ?`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
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
	return tx.GetContext(context.Background(), model, sql, values...)
}

func (b *BaseRepository) GetByIdTx(tx *sqlx.Tx, table string, model any, value string) error {
	return b.GetByColumnTx(tx, table, "id", value, model)
}

func (b *BaseRepository) SelectWhereTx(tx *sqlx.Tx, table string, model any, where map[string]any) error {
	placeholders, values := b.BuildWhere(where)

	if len(placeholders) > 0 && len(values) > 0 {
		sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", b.Escape(table), strings.Join(placeholders, " AND "))
		return tx.GetContext(context.Background(), model, sql, values...)

	}
	return errors.New("SelectWhere: Operation not allowed with empty where")
}

func (b *BaseRepository) DeleteByIdTx(tx *sqlx.Tx, table, id string) (sql.Result, error) {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id = ?`,
		b.Escape(table),
	)

	return tx.ExecContext(
		context.Background(),
		query,
		id,
	)
}

func (b *BaseRepository) InsertTx(tx *sqlx.Tx, table string, data map[string]any) (sql.Result, error) {

	placeholders, columns, values := b.BuildInsertPlaceHolders(data)

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		b.Escape(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return tx.ExecContext(
		context.Background(),
		query,
		values...,
	)
}

func (b *BaseRepository) InsertOrUpdateTx(tx *sqlx.Tx, table string, data map[string]any, where map[string]any) (sql.Result, error) {
	// check exists
	model := struct {
		Result string `db:"result"`
		ID     string `db:"id"`
	}{}

	err := b.SelectWhere(table, model, where)
	if err != nil {
		return b.InsertTx(tx, table, data)
	}
	return b.UpdateByIdTx(tx, table, model.ID, data)
}

func (b *BaseRepository) UpdateByIdTx(tx *sqlx.Tx, table string, id string, data map[string]any) (sql.Result, error) {
	placeholders, values := b.BuildUpdatePlaceHolders(data)
	query := fmt.Sprintf(
		`UPDATE %s SET %s  WHERE id = ?`,
		b.Escape(table),
		strings.Join(placeholders, ", "),
	)
	return tx.ExecContext(
		context.Background(),
		query,
		append([]any{id}, values...),
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
