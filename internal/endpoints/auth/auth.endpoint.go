package authendpoint

import (
	"context"
	appconfig "goapp/config"
	dbpostgres "goapp/internal/database/postgres"
	"goapp/internal/store"
	usersentity "goapp/internal/store/entities/users"
)

type Endpoint interface {
	Register(ctx context.Context, req RegisterRequest) RegisterResponse
}

type Auth struct {
	Config appconfig.Config
	Store  *store.Base
}

func New(cnf appconfig.Config, store *store.Base) Endpoint {
	return &Auth{
		Config: cnf,
		Store:  store,
	}
}

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
	ClientIP string
}

type RegisterResponse struct {
	UserId  int    `json:"user_id"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func (b Auth) Register(ctx context.Context, req RegisterRequest) RegisterResponse {

	metadata, _ := dbpostgres.ScanRawDataJSON(map[string]interface{}{
		"clientIp": req.ClientIP,
		"nextStep": "verify",
	})

	usr, err := b.Store.Postgres.User.Insert(ctx, usersentity.UserModels{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Status:   usersentity.UserStatusPending,
		MetaData: metadata,
	})
	if err != nil {
		return RegisterResponse{
			Status:  false,
			Message: err.Error(),
		}
	}
	return RegisterResponse{
		Status:  true,
		Message: "User created successfully",
		UserId:  usr.Id,
	}
}
