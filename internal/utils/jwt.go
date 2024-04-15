package utils

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret []byte

type JwtStandardClaims struct {
	IdentityKey string
	SecretValue string
	MetaData    map[string]string
	jwt.StandardClaims
}

type JwtStandardOptions struct {
	IdentityKey string
	SecretValue string
	MetaData    map[string]string
	Issuer      string
	Subject     string
	NoExpire    bool
	ExpireTime  time.Time
}

func JwtGenerateToken(r JwtStandardOptions) string {
	standardClaims := jwt.StandardClaims{
		Issuer:  r.Issuer,
		Subject: r.Subject,
	}
	if !r.NoExpire {
		standardClaims.ExpiresAt = r.ExpireTime.Unix()
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtStandardClaims{
		r.IdentityKey,
		r.SecretValue,
		r.MetaData,
		standardClaims,
	})
	value, err := claims.SignedString(jwtSecret)
	token := value
	if err != nil {
		token = ""
	}
	return token
}

func JwtParseToken(token string) *JwtStandardClaims {

	_Claims, claimsErr := jwt.ParseWithClaims(token, &JwtStandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claimsErr == nil && _Claims != nil {
		if claims, ok := _Claims.Claims.(*JwtStandardClaims); ok && _Claims.Valid {
			return claims
		}
	}
	return nil
}

func GenerateSha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}
