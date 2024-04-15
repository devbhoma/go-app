package authorization

import (
	"github.com/gin-gonic/gin"
	appconfig "goapp/config"
	"goapp/internal/store"
	"goapp/internal/utils"
	"net/http"
	"strings"
)

const (
	CookieKeyName     = "__app__user"
	CookieKeyMaxAge   = 60 * 60 * 24 * 365
	CookieKeySecure   = false
	CookieKeyHttpOnly = false
)

type Base struct {
	Config appconfig.Config
	Store  *store.Base
}

type Authorization interface {
	Authenticate() gin.HandlerFunc
	ValidateCookieToken(ctx *gin.Context) bool
}

func New(conf appconfig.Config, store *store.Base) Authorization {
	return &Base{
		Config: conf,
		Store:  store,
	}
}

func (a Base) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if a.ValidateCookieToken(ctx) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (a Base) ValidateCookieToken(ctx *gin.Context) bool {
	authorization := len(ctx.GetHeader("Authorization")) > 0
	cookie, _ := ctx.Cookie(CookieKeyName)

	if !authorization && cookie == "" {
		return false
	}

	token := ""
	if authorization {
		token = ctx.GetHeader("Authorization")
		if !strings.Contains(token, "Bearer ") {
			return false
		}
		token = strings.ReplaceAll(token, "Bearer", "")
		token = strings.TrimSpace(token)
	}
	if token == "" && cookie != "" {
		token = cookie
	}
	if token == "" {
		return false
	}

	parse := utils.JwtParseToken(token)
	if parse.SecretValue != "" {
		// db data validate
		ctx.Set("SECRET_KEY", parse.IdentityKey)
		ctx.Set("SECRET_VALUE", parse.SecretValue)
		return true
	}
	return false
}
