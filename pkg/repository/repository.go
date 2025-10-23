package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/sqlite"

	"github.com/K0ng2/zeedzad/db"
	"github.com/K0ng2/zeedzad/model"
)

type Repository struct {
	db *db.Database
	ex db.Executor
}

func NewRepository(db *db.Database) *Repository {
	return &Repository{db: db, ex: db.Conn()}
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{
		db: r.db,
		ex: tx,
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func FormatError(prefix string, err error) error {
	return fmt.Errorf("%s: %v", prefix, err)
}

func NullString(s string) sqlite.StringExpression {
	if strings.TrimSpace(s) == "" {
		return sqlite.StringExp(sqlite.NULL)
	}

	return sqlite.String(s)
}

func NullInt16(i *int) sqlite.IntegerExpression {
	if i == nil {
		return sqlite.IntExp(sqlite.NULL)
	}

	return sqlite.Int16(int16(*i))
}

func ILIKE(lhs, rhs sqlite.StringExpression) sqlite.BoolExpression {
	return sqlite.BoolExp(sqlite.BinaryOperator(lhs, rhs, "ILIKE"))
}

func TotalItems(ctx context.Context, exec db.Executor, countColumn sqlite.Expression, table sqlite.ReadableTable, expression *sqlite.BoolExpression) (int64, error) {
	var count model.INT64

	stmt := sqlite.SELECT(
		sqlite.COUNT(countColumn).AS("int64.number"),
	).FROM(table)

	if expression != nil {
		stmt = stmt.WHERE(*expression)
	}

	err := stmt.QueryContext(ctx, exec, &count)
	if err != nil {
		return 0, FormatError("creator", err)
	}

	return count.Number, nil
}
