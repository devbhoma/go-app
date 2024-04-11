package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	appconfig "goapp/config"
	"goapp/internal/authorization"
	authendpoint "goapp/internal/endpoints/auth"
	"goapp/internal/store"
	"goapp/internal/utils"
	"net/http"
)

type AuthHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type Endpoints struct {
	Auth authendpoint.Endpoint
}

type Auth struct {
	Store  *store.Base
	Config appconfig.Config
	Endpoints
}

func NewAuthHandler(cnf appconfig.Config, str *store.Base, router *gin.RouterGroup) AuthHandler {
	b := &Auth{
		Store:  str,
		Config: cnf,
		Endpoints: Endpoints{
			Auth: authendpoint.New(cnf, str),
		},
	}

	router.POST("/login", b.Login)
	router.POST("/register", b.Register)
	return b
}

func (a Auth) Login(ctx *gin.Context) {

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("req:", req)

	if req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is empty"})
		return
	}

	// connect with user entity

	// temporary fake token generated
	token := utils.JwtGenerateToken(utils.JwtStandardOptions{
		IdentityKey: "static_user",
		SecretValue: "static_user_fake_id",
		MetaData: map[string]string{
			"ip": ctx.ClientIP(),
		},
		NoExpire: true,
	})

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(authorization.CookieKeyName, token, authorization.CookieKeyMaxAge, "/", "", authorization.CookieKeySecure, authorization.CookieKeyHttpOnly)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a Auth) Register(ctx *gin.Context) {

	type Request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name or email or password is empty"})
		return
	}

	resp := a.Endpoints.Auth.Register(ctx, authendpoint.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		ClientIP: ctx.ClientIP(),
	})

	ctx.JSON(http.StatusOK, resp)
}
