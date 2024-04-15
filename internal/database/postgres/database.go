package dbpostgres

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sqlc-dev/pqtype"
)

type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

func StatusCheck(ctx context.Context, db *sql.DB) error {
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

	return &DB{SqlDB: db}, nil
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

func (db *DB) Ping() error {
	return db.SqlDB.Ping()
}


func (db *DB) Close() error {
	return db.SqlDB.Close()
}

func (db *DB) RawSqlDB() *sql.DB {
	return db.SqlDB
}


func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.SqlDB.QueryRowContext(ctx, query, args...)
}


func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.SqlDB.QueryContext(ctx, query, args...)
}
func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.SqlDB.ExecContext(ctx, query, args...)
}
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.SqlDB.BeginTx(ctx, opts)
}
