package store

import (
	"context"
	"fmt"
	appconfig "goapp/config"
	dbpostgres "goapp/internal/database/postgres"
	"goapp/internal/database/redis"
	usersentity "goapp/internal/store/entities/users"
	migrationentity "goapp/internal/store/migration"
)

type RedisStore struct {
	Reader redis.Reader
}

type PostgresStore struct {
	DB        *dbpostgres.DB
	User      usersentity.User
	Migration migrationentity.Migration
}
type Base struct {
	Postgres PostgresStore
	Redis    RedisStore
}

func NewStore(cnf appconfig.Config) *Base {

	return &Base{
		Postgres: NewPgDb(cnf),
		Redis:    RedisStore{},
	}
}

func NewPgDb(cnf appconfig.Config) PostgresStore {

	db, err := dbpostgres.New(
		cnf.Database.Driver,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.User,
		cnf.Database.Pass,
		cnf.Database.Name,
	)
	if err != nil {
		panic(err)
	}

	if err := dbpostgres.StatusCheck(context.Background(), db.RawSqlDB()); err != nil {
		fmt.Printf("cannot connect to postgres db", "err", err)
		panic(err)
	}
	return NewPgDbStore(db, cnf)

}

func NewPgDbStore(db *dbpostgres.DB, cnf appconfig.Config) PostgresStore {
	return PostgresStore{
		DB:        db,
		User:      usersentity.NewUserStore(db),
		Migration: migrationentity.NewStore(db),
	}
}
