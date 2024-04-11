package migrationentity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	dbpostgres "goapp/internal/database/postgres"
	"time"
)

type Migration interface {
	Insert(context.Context, MigrationModels) (MigrationModels, error)
	GetByName(context.Context, string) (MigrationModels, error)
	GetAll(context.Context) ([]MigrationModels, error)
	Create(context.Context, []map[string]interface{}) (bool, error)
	CreateOneByOne(context.Context, string) bool
}

type MigrationModels struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MigrationBase struct {
	Database *dbpostgres.DB
}

func NewStore(db *dbpostgres.DB) Migration {
	return &MigrationBase{
		Database: db,
	}
}

func (c *MigrationBase) Insert(ctx context.Context, arg MigrationModels) (MigrationModels, error) {
	var r MigrationModels
	qry := c.Database.QueryRow(ctx,
		`INSERT INTO public.db_migrations( name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, name, created_at, updated_at`,
		arg.Name, arg.CreatedAt, arg.UpdatedAt)
	err := qry.Scan(&r.Id, &r.Name, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (c *MigrationBase) GetByName(ctx context.Context, name string) (MigrationModels, error) {
	var r MigrationModels
	qry := c.Database.QueryRow(ctx, `SELECT * FROM public.db_migrations WHERE name = $1`, name)
	err := qry.Scan(&r.Id, &r.Name, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return r, err
}

func (c *MigrationBase) GetAll(ctx context.Context) ([]MigrationModels, error) {
	var r []MigrationModels
	qry, err := c.Database.Query(ctx, `SELECT * FROM public.db_migrations`)
	if err == nil {
		for qry.Next() {
			var m MigrationModels
			err = qry.Scan(&m.Id, &m.Name, &m.CreatedAt, &m.UpdatedAt)
			if err == nil {
				r = append(r, m)
			}
		}
	}
	return r, err
}

func (c *MigrationBase) Create(ctx context.Context, queries []map[string]interface{}) (bool, error) {
	tx, err := c.Database.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	for _, val := range queries {
		_, err = tx.ExecContext(ctx, val["qry"].(string))
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				return false, rErr
			}
		}
	}
	err1 := tx.Commit()
	return err1 == nil, err1
}

func (c *MigrationBase) CreateOneByOne(ctx context.Context, qry string) bool {
	tx, bErr := c.Database.BeginTx(ctx, nil)
	if bErr == nil {
		_, eErr := tx.ExecContext(ctx, qry)
		if eErr != nil {
			fmt.Println("ExecContext Err: ", eErr)
			if rErr := tx.Rollback(); rErr != nil {
				fmt.Println("Rollback Err: ", eErr)
			}
			return false
		}

		if err1 := tx.Commit(); err1 != nil {
			fmt.Println("Commit Err: ", eErr)
			return false
		}
		return true
	}
	fmt.Println("BeginTx Err: ", bErr)
	return false
}
