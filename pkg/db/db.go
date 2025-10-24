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

type Database struct {
	client     *cloudflare.Client
	accountID  string
	databaseID string
}

func (d *Database) Close() error {
	// No-op for HTTP client
	return nil
}

func (d *Database) Conn() *sql.DB {
	// Return a database/sql compatible connection
	// Since go-jet requires this, we need to provide a wrapper
	return d.getOrCreateSQLDB()
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// D1 doesn't support traditional transactions via REST API
	// We'll simulate it by returning a fake tx that executes immediately
	return nil, fmt.Errorf("transactions not supported with D1 REST API")
}

func (d *Database) PingContext(ctx context.Context) error {
	// Test the connection with a simple query
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

// inlineParams replaces ? placeholders with actual values for batch queries
// This is needed because D1 doesn't support params with multiple statements
func inlineParams(query string, args ...any) string {
	result := query
	for _, arg := range args {
		// Find the first ? and replace it
		idx := strings.Index(result, "?")
		if idx == -1 {
			break
		}

		var replacement string
		if arg == nil {
			replacement = "NULL"
		} else {
			// Format the argument based on type
			switch v := arg.(type) {
			case string:
				// Escape single quotes in strings
				escaped := strings.ReplaceAll(v, "'", "''")
				replacement = "'" + escaped + "'"
			case *string:
				if v == nil {
					replacement = "NULL"
				} else {
					escaped := strings.ReplaceAll(*v, "'", "''")
					replacement = "'" + escaped + "'"
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				replacement = fmt.Sprintf("%d", v)
			case float32, float64:
				replacement = fmt.Sprintf("%f", v)
			case bool:
				if v {
					replacement = "1"
				} else {
					replacement = "0"
				}
			case time.Time:
				// Format time as SQLite datetime format
				replacement = "'" + v.Format("2006-01-02 15:04:05") + "'"
			default:
				// For other types, convert to string and quote
				str := fmt.Sprintf("%v", v)
				escaped := strings.ReplaceAll(str, "'", "''")
				replacement = "'" + escaped + "'"
			}
		}

		// Replace the first ? with the value
		result = result[:idx] + replacement + result[idx+1:]
	}
	return result
}

// Execute raw SQL query via D1 API
func (d *Database) executeQuery(ctx context.Context, query string, args ...any) (*QueryResult, error) {
	// Check if this is a write operation (INSERT, UPDATE, DELETE)
	// Defer foreign key constraints for write operations
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isWriteOp := strings.HasPrefix(upperQuery, "INSERT") ||
		strings.HasPrefix(upperQuery, "UPDATE") ||
		strings.HasPrefix(upperQuery, "DELETE")

	var finalQuery string
	var params []string

	if isWriteOp {
		// For write operations with PRAGMA, we need to inline parameters
		// because D1 doesn't support params with multiple statements
		finalQuery = "PRAGMA defer_foreign_keys = on; " + inlineParams(query, args...)
		params = nil // No separate params when using batch queries
	} else {
		// For read operations, use parameterized queries
		finalQuery = query
		params = make([]string, len(args))
		for i, arg := range args {
			if arg == nil {
				params[i] = "NULL"
			} else {
				params[i] = fmt.Sprintf("%v", arg)
			}
		}
	}

	resp, err := d.client.D1.Database.Query(ctx, d.databaseID, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(d.accountID),
		Sql:       cloudflare.F(finalQuery),
		Params:    cloudflare.F(params),
	})
	if err != nil {
		return nil, err
	}

	// Get the result (D1 returns an array of results for batch queries)
	if len(resp.Result) == 0 {
		return nil, fmt.Errorf("no results returned from D1")
	}

	// If we prepended PRAGMA, we'll have 2 results: [PRAGMA result, actual query result]
	// Use the last result which is the actual query
	d1Result := resp.Result[len(resp.Result)-1]

	// Convert results from []interface{} to []map[string]interface{}
	results := make([]map[string]interface{}, 0, len(d1Result.Results))
	for _, r := range d1Result.Results {
		if m, ok := r.(map[string]interface{}); ok {
			results = append(results, m)
		}
	}

	result := &QueryResult{
		Success: d1Result.Success,
		Meta: QueryResultMeta{
			ChangedDB:   d1Result.Meta.ChangedDB,
			Changes:     d1Result.Meta.Changes,
			Duration:    d1Result.Meta.Duration,
			LastRowID:   d1Result.Meta.LastRowID,
			RowsRead:    d1Result.Meta.RowsRead,
			RowsWritten: d1Result.Meta.RowsWritten,
			SizeAfter:   d1Result.Meta.SizeAfter,
		},
		Results: results,
	}

	if !result.Success {
		return nil, fmt.Errorf("D1 query failed")
	}

	return result, nil
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
	queryArgs := make([]any, len(args))
	for i, arg := range args {
		queryArgs[i] = arg.Value
	}

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
	queryArgs := make([]any, len(args))
	for i, arg := range args {
		queryArgs[i] = arg.Value
	}

	result, err := c.database.executeQuery(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return &d1Rows{results: result.Results, index: -1}, nil
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
	return s.ExecContext(context.Background(), toNamedValues(args))
}

func (s *d1Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(context.Background(), toNamedValues(args))
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
		r.columns = make([]string, 0, len(r.results[0]))
		for col := range r.results[0] {
			r.columns = append(r.columns, col)
		}
	}
	return r.columns
}

func (r *d1Rows) Close() error {
	return nil
}

func (r *d1Rows) Next(dest []driver.Value) error {
	r.index++
	if r.index >= len(r.results) {
		return io.EOF
	}

	row := r.results[r.index]
	columns := r.Columns()

	for i, col := range columns {
		if val, ok := row[col]; ok {
			// Convert datetime strings to a format go-jet can handle
			if strVal, isString := val.(string); isString {
				// Try various datetime formats
				var t time.Time
				var err error

				// Try to parse datetime - strip monotonic clock if present
				// Handle format: "2025-10-24 16:48:30.211971376 +0700 +07 m=+50.122713380"
				parseStr := strVal
				if idx := strings.Index(parseStr, " m="); idx > 0 {
					parseStr = parseStr[:idx] // Remove monotonic clock portion
				}

				// Try various datetime formats
				// RFC3339: "2025-10-22T09:00:27Z"
				if t, err = time.Parse(time.RFC3339, parseStr); err == nil {
					dest[i] = t.Format("2006-01-02 15:04:05")
				} else if t, err = time.Parse(time.RFC3339Nano, parseStr); err == nil {
					dest[i] = t.Format("2006-01-02 15:04:05")
				} else if t, err = time.Parse("2006-01-02 15:04:05 -0700 MST", parseStr); err == nil {
					// Go time.String() format: "2025-03-05 09:00:39 +0000 UTC"
					dest[i] = t.Format("2006-01-02 15:04:05")
				} else if t, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", parseStr); err == nil {
					// Go time.String() format with nanoseconds
					dest[i] = t.Format("2006-01-02 15:04:05")
				} else if t, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 -07", parseStr); err == nil {
					// Go time.String() format with numeric timezone
					dest[i] = t.Format("2006-01-02 15:04:05")
				} else if _, err = time.Parse("2006-01-02 15:04:05", parseStr); err == nil {
					// Already in SQLite format
					dest[i] = parseStr
				} else if _, err = time.Parse("2006-01-02", parseStr); err == nil {
					// Date only
					dest[i] = parseStr
				} else {
					// Not a datetime, keep as string
					dest[i] = val
				}
			} else {
				dest[i] = val
			}
		} else {
			dest[i] = nil
		}
	}

	return nil
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

func toNamedValues(values []driver.Value) []driver.NamedValue {
	named := make([]driver.NamedValue, len(values))
	for i, v := range values {
		named[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return named
}

func NewDatabase(accountID, databaseID, apiToken string) (*Database, error) {
	if accountID == "" || databaseID == "" || apiToken == "" {
		return nil, fmt.Errorf("accountID, databaseID, and apiToken are required")
	}

	client := cloudflare.NewClient(option.WithAPIToken(apiToken))

	return &Database{
		client:     client,
		accountID:  accountID,
		databaseID: databaseID,
	}, nil
}
