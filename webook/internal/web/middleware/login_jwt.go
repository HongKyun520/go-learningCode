package middleware

import (
	"GoInAction/webook/internal/web"
	"encoding/gob"
	"net/http"
	"time"

	ijwt "GoInAction/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 基于jwt实现的登录校验
type LoginJWTMiddlewareBuilder struct {
	ignorePaths []string
	handler     ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(handler ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		handler: handler,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		// 不需要登录校验的
		// 不需要登录状态校验的
		// 从忽略中的paths取出遍历

		for _, path := range l.ignorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 获取jwt
		tokenStr := web.ExtractToken(ctx)

		claims := &web.UserClaims{}

		// 校验jwt
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"), nil
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题，需要打印记录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 设置uid
		ctx.Set("userId", claims.Uid)
	}
}
