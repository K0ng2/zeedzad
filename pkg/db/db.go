package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/d1"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

const (
	sqliteDateTimeFormat = "2006-01-02 15:04:05"
	sqliteDateFormat     = "2006-01-02"
)

type Database struct {
	client     *cloudflare.Client
	accountID  string
	databaseID string
}

func NewDatabase(accountID, databaseID, apiToken string) (*Database, error) {
	if err := validateDatabaseConfig(accountID, databaseID, apiToken); err != nil {
		return nil, err
	}

	client := cloudflare.NewClient(option.WithAPIToken(apiToken))

	return &Database{
		client:     client,
		accountID:  accountID,
		databaseID: databaseID,
	}, nil
}

func validateDatabaseConfig(accountID, databaseID, apiToken string) error {
	if accountID == "" || databaseID == "" || apiToken == "" {
		return fmt.Errorf("accountID, databaseID, and apiToken are required")
	}
	return nil
}

func (d *Database) Close() error {
	return nil
}

func (d *Database) Conn() *sql.DB {
	return d.getOrCreateSQLDB()
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, fmt.Errorf("transactions not supported with D1 REST API")
}

func (d *Database) PingContext(ctx context.Context) error {
	_, err := d.client.D1.Database.Query(ctx, d.databaseID, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(d.accountID),
		Sql:       cloudflare.F("SELECT 1"),
	})
	return err
}

// D1 Result structures
type QueryResult struct {
	Success bool                     `json:"success"`
	Meta    QueryResultMeta          `json:"meta"`
	Results []map[string]interface{} `json:"results"`
}

type QueryResultMeta struct {
	ChangedDB   bool    `json:"changed_db"`
	Changes     float64 `json:"changes"`
	Duration    float64 `json:"duration"`
	LastRowID   float64 `json:"last_row_id"`
	RowsRead    float64 `json:"rows_read"`
	RowsWritten float64 `json:"rows_written"`
	SizeAfter   float64 `json:"size_after"`
}

// D1 Executor implementation
type Executor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// D1Result implements sql.Result interface
type D1Result struct {
	lastInsertID int64
	rowsAffected int64
}

func (r D1Result) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

func (r D1Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func inlineParams(query string, args ...any) string {
	result := query
	for _, arg := range args {
		idx := strings.Index(result, "?")
		if idx == -1 {
			break
		}

		replacement := formatInlineParam(arg)
		result = result[:idx] + replacement + result[idx+1:]
	}
	return result
}

func formatInlineParam(arg any) string {
	if arg == nil {
		return "NULL"
	}

	switch v := arg.(type) {
	case string:
		return formatStringParam(v)
	case *string:
		return formatStringPointerParam(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return formatBoolParam(v)
	case time.Time:
		return formatTimeParam(v)
	default:
		return formatDefaultParam(v)
	}
}

func formatStringParam(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}

func formatStringPointerParam(s *string) string {
	if s == nil {
		return "NULL"
	}
	return formatStringParam(*s)
}

func formatBoolParam(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func formatTimeParam(t time.Time) string {
	return "'" + t.Format(sqliteDateTimeFormat) + "'"
}

func formatDefaultParam(v any) string {
	str := fmt.Sprintf("%v", v)
	escaped := strings.ReplaceAll(str, "'", "''")
	return "'" + escaped + "'"
}

func (d *Database) executeQuery(ctx context.Context, query string, args ...any) (*QueryResult, error) {
	finalQuery, params := d.prepareQuery(query, args...)

	resp, err := d.client.D1.Database.Query(ctx, d.databaseID, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(d.accountID),
		Sql:       cloudflare.F(finalQuery),
		Params:    cloudflare.F(params),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Result) == 0 {
		return nil, fmt.Errorf("no results returned from D1")
	}

	d1Result := resp.Result[len(resp.Result)-1]

	if !d1Result.Success {
		return nil, fmt.Errorf("D1 query failed")
	}

	results := convertD1Results(d1Result.Results)
	meta := convertD1Meta(d1Result.Meta)

	return &QueryResult{
		Success: d1Result.Success,
		Meta:    meta,
		Results: results,
	}, nil
}

func (d *Database) prepareQuery(query string, args ...any) (string, []string) {
	if isWriteOperation(query) {
		return d.prepareWriteQuery(query, args...)
	}
	return d.prepareReadQuery(query, args...)
}

func isWriteOperation(query string) bool {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(upperQuery, "INSERT") ||
		strings.HasPrefix(upperQuery, "UPDATE") ||
		strings.HasPrefix(upperQuery, "DELETE")
}

func (d *Database) prepareWriteQuery(query string, args ...any) (string, []string) {
	inlinedQuery := "PRAGMA defer_foreign_keys = on; " + inlineParams(query, args...)
	return inlinedQuery, nil
}

func (d *Database) prepareReadQuery(query string, args ...any) (string, []string) {
	params := make([]string, len(args))
	for i, arg := range args {
		params[i] = formatParamValue(arg)
	}
	return query, params
}

func formatParamValue(arg any) string {
	if arg == nil {
		return "NULL"
	}
	return fmt.Sprintf("%v", arg)
}

func convertD1Results(d1Results []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(d1Results))
	for _, r := range d1Results {
		if m, ok := r.(map[string]interface{}); ok {
			results = append(results, m)
		}
	}
	return results
}

func convertD1Meta(meta interface{}) QueryResultMeta {
	type metaStruct struct {
		ChangedDB   bool
		Changes     float64
		Duration    float64
		LastRowID   float64
		RowsRead    float64
		RowsWritten float64
		SizeAfter   float64
	}

	if m, ok := meta.(metaStruct); ok {
		return QueryResultMeta{
			ChangedDB:   m.ChangedDB,
			Changes:     m.Changes,
			Duration:    m.Duration,
			LastRowID:   m.LastRowID,
			RowsRead:    m.RowsRead,
			RowsWritten: m.RowsWritten,
			SizeAfter:   m.SizeAfter,
		}
	}

	return QueryResultMeta{}
}

// getOrCreateSQLDB creates a fake database/sql.DB that wraps D1 API
func (d *Database) getOrCreateSQLDB() *sql.DB {
	// Register a custom driver for D1
	driverName := fmt.Sprintf("d1-%s-%s", d.accountID, d.databaseID)

	sql.Register(driverName, &d1Driver{
		database: d,
	})

	db, _ := sql.Open(driverName, "")
	return db
}

// D1 driver implementation
type d1Driver struct {
	database *Database
}

func (drv *d1Driver) Open(name string) (driver.Conn, error) {
	return &d1Conn{database: drv.database}, nil
}

type d1Conn struct {
	database *Database
}

func (c *d1Conn) Prepare(query string) (driver.Stmt, error) {
	return &d1Stmt{conn: c, query: query}, nil
}

func (c *d1Conn) Close() error {
	return nil
}

func (c *d1Conn) Begin() (driver.Tx, error) {
	return &d1Tx{conn: c}, nil
}

func (c *d1Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	queryArgs := convertNamedValuesToArgs(args)

	result, err := c.database.executeQuery(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return D1Result{
		lastInsertID: int64(result.Meta.LastRowID),
		rowsAffected: int64(result.Meta.Changes),
	}, nil
}

func (c *d1Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryArgs := convertNamedValuesToArgs(args)

	result, err := c.database.executeQuery(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &d1Rows{results: result.Results, index: -1}, nil
}

func convertNamedValuesToArgs(args []driver.NamedValue) []any {
	queryArgs := make([]any, len(args))
	for i, arg := range args {
		queryArgs[i] = arg.Value
	}
	return queryArgs
}

type d1Stmt struct {
	conn  *d1Conn
	query string
}

func (s *d1Stmt) Close() error {
	return nil
}

func (s *d1Stmt) NumInput() int {
	return strings.Count(s.query, "?")
}

func (s *d1Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.ExecContext(context.Background(), convertValuesToNamedValues(args))
}

func (s *d1Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(context.Background(), convertValuesToNamedValues(args))
}

func (s *d1Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.conn.ExecContext(ctx, s.query, args)
}

func (s *d1Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return s.conn.QueryContext(ctx, s.query, args)
}

type d1Rows struct {
	results []map[string]interface{}
	index   int
	columns []string
}

func (r *d1Rows) Columns() []string {
	if r.columns == nil && len(r.results) > 0 {
		r.columns = extractColumnNames(r.results[0])
	}
	return r.columns
}

func extractColumnNames(row map[string]interface{}) []string {
	columns := make([]string, 0, len(row))
	for col := range row {
		columns = append(columns, col)
	}
	return columns
}

func (r *d1Rows) Close() error {
	return nil
}

func (r *d1Rows) Next(dest []driver.Value) error {
	r.index++
	if r.index >= len(r.results) {
		return io.EOF
	}

	r.populateRowValues(dest)
	return nil
}

func (r *d1Rows) populateRowValues(dest []driver.Value) {
	row := r.results[r.index]
	columns := r.Columns()

	for i, col := range columns {
		if val, ok := row[col]; ok {
			dest[i] = convertColumnValue(val)
		} else {
			dest[i] = nil
		}
	}
}

func convertColumnValue(val interface{}) interface{} {
	strVal, isString := val.(string)
	if !isString {
		return val
	}

	return parseDateTime(strVal)
}

func parseDateTime(strVal string) interface{} {
	cleanedStr := stripMonotonicClock(strVal)

	timeFormats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999999 -0700 -07",
		sqliteDateTimeFormat,
		sqliteDateFormat,
	}

	for _, format := range timeFormats {
		if t, err := time.Parse(format, cleanedStr); err == nil {
			if format == sqliteDateTimeFormat || format == sqliteDateFormat {
				return cleanedStr
			}
			return t.Format(sqliteDateTimeFormat)
		}
	}

	return strVal
}

func stripMonotonicClock(s string) string {
	if idx := strings.Index(s, " m="); idx > 0 {
		return s[:idx]
	}
	return s
}

type d1Tx struct {
	conn *d1Conn
}

func (tx *d1Tx) Commit() error {
	// No-op for D1
	return nil
}

func (tx *d1Tx) Rollback() error {
	// No-op for D1
	return nil
}

func convertValuesToNamedValues(values []driver.Value) []driver.NamedValue {
	named := make([]driver.NamedValue, len(values))
	for i, v := range values {
		named[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return named
}
