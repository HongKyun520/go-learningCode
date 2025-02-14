package ioc

import (
	"GoInAction/webook/internal/web"
	"GoInAction/webook/internal/web/middleware"

	// "GoInAction/webook/internal/web/middleware"
	"GoInAction/webook/pkg/ginx/middlewares/prometheus"
	"GoInAction/webook/pkg/logger"

	// "GoInAction/webook/pkg/ginx/middlewares/ratelimit"
	// "GoInAction/webook/pkg/limiter"
	ijwt "GoInAction/webook/internal/web/jwt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	oauth2Hdl *web.OAuth2WechatHandler,
	artHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterUsersRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	return server
}

// 初始化中间件
func InitMiddlewares(redisClient redis.Cmdable, l logger.Logger, handler ijwt.Handler) []gin.HandlerFunc {

	pb := &prometheus.Builder{
		NameSpace:  "geektime_daming",
		Subsystem:  "webook",
		Name:       "gin_http",
		InstanceId: "1234567890",
	}

	return []gin.HandlerFunc{

		// 解决跨域问题
		cors.New(cors.Config{
			// 允许的请求来源（一般填写域名）
			//AllowOrigins: []string{"https://localhost:8080"},
			// 允许的请求方法 不写默认都支持

			//AllowMethods: []string{"PUT", "PATCH", "POST"},

			// 允许的请求携带的请求头
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},

			// 允许暴露的响应头
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},

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

		// 打印请求日志
		func(ctx *gin.Context) {
			println("这是我的middleware")
		},

		// 限流
		// ratelimit.NewBuilder(limiter.NewRedisSlideWindowLimiter(redisClient, time.Minute, 1000)).Build(),

		// 登录校验
		// (&middleware.LoginJWTMiddlewareBuilder{}).
		// 	IgnorePaths("/users/loginJWT").
		// 	IgnorePaths("/users/signup").
		// 	IgnorePaths("/users/login").
		// 	IgnorePaths("/users/login_sms/code/send").
		// 	IgnorePaths("/users/login_sms").
		// 	IgnorePaths("/oauth2/wechat/authurl").
		// 	IgnorePaths("/oauth2/wechat/callback").
		// 	Build(),
		pb.BuildResponseTime(),
		pb.BuildActiveRequests(),

		middleware.NewLogMiddlewareBuilder(func(ctx *gin.Context, al middleware.AccessLog) {
			l.Debug("access log", logger.Field{
				Key:   "request",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Build(),

		middleware.NewLoginJWTMiddlewareBuilder(handler).Build(),
	}
}
