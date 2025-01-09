package ioc

import (
	"GoInAction/webook/internal/web"
	"GoInAction/webook/internal/web/middleware"
	"GoInAction/webook/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterUsersRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{

			// 允许的请求来源（一般填写域名）
			//AllowOrigins: []string{"https://localhost:8080"},
			// 允许的请求方法 不写默认都支持

			//AllowMethods: []string{"PUT", "PATCH", "POST"},

			// 允许的请求携带的请求头
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},

			// 允许暴露的响应头
			ExposeHeaders: []string{"x-jwt-token"},

			// 是否允许携带cookie之类的东西
			AllowCredentials: true,

			// 自定义路由策略，使用方法进行扩展
			AllowOriginFunc: func(origin string) bool {
				// 测试环境
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}

				return strings.Contains(origin, "your domain")
			},
			// preflight请求有效期
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("这是我的middleware")
		},
		ratelimit.NewBuilder(redisClient, time.Minute, 1000).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).
			IgnorePaths("/users/loginJWT").
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			Build(),
	}
}
