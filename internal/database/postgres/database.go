package dbpostgres

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sqlc-dev/pqtype"
	"net/url"
)

type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

func Open(cfg Config) (*sql.DB, error) {

	// Define SSL mode.
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sql.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

func New(driver, host, port, user, pwd, dbname string) (*DB, error) {
	var source string
	source = "port=" + port
	source += " host=" + host
	source += " user=" + user
	if pwd != "" {
		source += " password=" + pwd
	}
	source += " dbname=" + dbname
	source += " sslmode=disable"

	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to primary db")
	}

	db.SetMaxOpenConns(400)
	db.SetMaxIdleConns(100)

	return &DB{SqlDB: db, Mapper: mapper()}, nil
}

func ScanRawDataJSON(v interface{}) (pqtype.NullRawMessage, error) {
	rawMessage := pqtype.NullRawMessage{}
	rawValueBuf, _ := json.Marshal(v)
	if err := rawMessage.Scan(rawValueBuf); err != nil {
		return rawMessage, err
	}
	return rawMessage, nil
}

func ScanRawDataString(v interface{}) (sql.NullString, error) {
	rawMessage := sql.NullString{}

	if err := rawMessage.Scan(v); err != nil {
		return rawMessage, err
	}
	return rawMessage, nil
}

func BindRawDataJSON(rawMessage pqtype.NullRawMessage, message interface{}) error {
	if rawMessage.Valid {
		var err error
		jsonKeys, err := json.Marshal(rawMessage.RawMessage)
		if err != nil {
			return err
		}
		_ = json.Unmarshal(jsonKeys, &message)
	}
	return nil
}

type DB struct {
	SqlDB  *sql.DB
	unsafe bool
	Mapper *Mapper
}

type DBQuery interface {
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	RawSqlDB() *sql.DB
	Ping() error
	Close() error
	GetContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error
}

// Ping pings the database to check the connection status
func (db *DB) Ping() error {
	return db.SqlDB.Ping()
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (db *DB) Close() error {
	return db.SqlDB.Close()
}

// RawSqlDB returns the sql.db implementions
// useful for anyone to execute queries outside of endpoints
func (db *DB) RawSqlDB() *sql.DB {
	return db.SqlDB
}

// Usually used to retrive only one row
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.SqlDB.QueryRowContext(ctx, query, args...)
}

// Used to retrive multiple rows in the table.
func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.SqlDB.QueryContext(ctx, query, args...)
}
func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.SqlDB.ExecContext(ctx, query, args...)
}
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.SqlDB.BeginTx(ctx, opts)
}

// GetContext using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) GetContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error {
	return GetContext(ctx, db.SqlDB, dest, q, args...)
}

// GetContext does a QueryRow using the provided Queryer, and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return sql.ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func GetContext(ctx context.Context, db *sql.DB, dest interface{}, query string, args ...interface{}) error {
	r := QueryRowContext(ctx, db, query, args...)
	return r.scanAny(dest, false)
}

// QueryRowContext queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func QueryRowContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) *Row {
	rows, err := db.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, unsafe: false, Mapper: mapper()}
}

// SelectContext using this DB.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return SelectContext(ctx, db.SqlDB, dest, query, args...)
}

// SelectContext executes a query using the provided Queryer, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func SelectContext(ctx context.Context, db *sql.DB, dest interface{}, query string, args ...interface{}) error {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsi := &Rows{Rows: rows, unsafe: false, Mapper: mapper()}
	defer rowsi.Close()
	return scanAll(rowsi, dest, false)
}
