package usersentity

import (
	"context"
	"fmt"
	"github.com/sqlc-dev/pqtype"
	dbpostgres "goapp/internal/database/postgres"
	"strings"
	"time"
)

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

type User interface {
	Insert(ctx context.Context, model UserModels) (UserModels, error)
	Get(ctx context.Context, opts UserOptions) (UserModels, error)
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

type UserOptions struct {
	Id       int
	Email    string
	Password string
	Status   string
}

func (u UserStore) Get(ctx context.Context, opts UserOptions) (UserModels, error) {
	qry := "SELECT id, name, email, status, meta_data FROM public.user "

	values := make([]interface{}, 0)
	where := make([]string, 0)

	if opts.Email != "" {
		values = append(values, opts.Email)
		where = append(where, fmt.Sprintf("email = $%d", len(values)))
	}
	if opts.Password != "" {
		values = append(values, opts.Password)
		where = append(where, fmt.Sprintf("password = $%d", len(values)))
	}
	if opts.Status != "" {
		values = append(values, opts.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(values)))
	}
	if opts.Id > 0 {
		where = append(where, fmt.Sprintf("id = $%d", len(values)))
		values = append(values, opts.Id)
	}

	if len(where) == 0 {
		return UserModels{}, nil
	}

	qry += " WHERE " + strings.Join(where, " AND ")

	row := u.Database.QueryRow(ctx, qry, values...)
	var usr UserModels
	err := row.Scan(&usr.Id, &usr.Name, &usr.Email, &usr.Status, &usr.MetaData)
	if err != nil {
		return UserModels{}, err
	}
	return usr, nil
}
