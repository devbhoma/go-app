package usersentity

import (
	"context"
	"github.com/sqlc-dev/pqtype"
	dbpostgres "goapp/internal/database/postgres"
	"time"
)

type User interface {
	Insert(ctx context.Context, model UserModels) (UserModels, error)
}

type UserStore struct {
	Database *dbpostgres.DB
}

type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
)

type UserModels struct {
	Id        int
	Name      string
	Email     string
	Password  string
	Status    UserStatus
	MetaData  pqtype.NullRawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserStore(db *dbpostgres.DB) User {
	return &UserStore{
		Database: db,
	}
}

func (u UserStore) Insert(ctx context.Context, arg UserModels) (UserModels, error) {

	arg.CreatedAt = time.Now()
	arg.UpdatedAt = time.Now()

	row := u.Database.QueryRow(ctx,
		`INSERT INTO public.user (name, email, password, status, meta_data, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		arg.Name, arg.Email, arg.Password, arg.Status, arg.MetaData, arg.CreatedAt, arg.UpdatedAt)

	var usr UserModels
	err := row.Scan(&usr.Id, &usr.Name, &usr.Email, &usr.Password, &usr.Status, &usr.MetaData, &usr.CreatedAt, &usr.UpdatedAt)
	usr.Password = ""
	return usr, err
}
