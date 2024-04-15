package authendpoint

import (
	"context"
	appconfig "goapp/config"
	dbpostgres "goapp/internal/database/postgres"
	"goapp/internal/store"
	usersentity "goapp/internal/store/entities/users"
	"goapp/internal/utils"
)

type Endpoint interface {
	Register(ctx context.Context, req RegisterRequest) RegisterResponse
	Login(ctx context.Context, req LoginRequest) LoginResponse
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	User    usersentity.UserModels
}

func (b Auth) Login(ctx context.Context, req LoginRequest) LoginResponse {

	usr, err := b.Store.Postgres.User.Get(ctx, usersentity.UserOptions{
		Email:    req.Username,
		Password: utils.GenerateSha1(req.Password),
	})

	if err != nil {
		return LoginResponse{
			Status:  false,
			Message: err.Error(),
		}
	}

	if usr.Id > 0 {
		return LoginResponse{
			Status:  true,
			Message: "User login successfully",
			User:    usr,
		}
	}

	return LoginResponse{
		Status:  false,
		Message: "credentials are invalid",
	}
}

func (b Auth) Register(ctx context.Context, req RegisterRequest) RegisterResponse {

	metadata, _ := dbpostgres.ScanRawDataJSON(map[string]interface{}{
		"clientIp": req.ClientIP,
		"nextStep": "verify",
	})

	usr, err := b.Store.Postgres.User.Insert(ctx, usersentity.UserModels{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.GenerateSha1(req.Password),
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
		UserId:  usr.Id,
		Message: "You have successfully registered",
	}
}
