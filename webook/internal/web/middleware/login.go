package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 基于session实现的登录校验

// 扩展性
type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
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

		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新token
		updateTime := sess.Get("update_time")
		sess.Set("userId", userId)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})

		now := time.Now().UnixMilli()
		if nil == updateTime {
			sess.Set("update_time", now)
		}

		updateTimeVal := updateTime.(int64)

		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
		}

	}
}
