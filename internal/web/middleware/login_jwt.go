package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-server/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuild() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// determine if the path need to checked
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}

		// is there a token
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claim := &web.UserClaim{}
		t, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !t.Valid || t == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claim.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
