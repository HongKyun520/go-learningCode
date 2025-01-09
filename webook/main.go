package main

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/service/sms/memory"
	"GoInAction/webook/internal/web"
	"GoInAction/webook/internal/web/middleware"
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// main方法需要自己手写很多东西

// go项目需要在main方法中手动做很多事情
// 初始化数据源、甚至初始化普通类、手动注册接口路由、手动注册middleware(过滤器、拦截器)

// 应用启动入口
func main() {

	//初始化服务器（使用wire）
	server := InitWebServer()

	//初始化数据库
	// db := initDB()

	//初始化用户控制器 controller -> service -> repository -> dao
	// u := initUser(db)

	// initRedis(server)

	//将userHandler里面的路由注册进server
	// u.RegisterUsersRoutes(server)
	server.Run(":8080")

	//engine := gin.Default()
	//engine.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello world")
	//})
	//
	//engine.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// middleware相当于spring中的拦截器，起到一个请求前置处理的作用
	//server.Use(func(context *gin.Context) {
	//	println("这是第一个middleware")
	//}, func(context *gin.Context) {
	//	println("这是第二个middleware")
	//})

	//// 初始化redis客户端
	//cmd := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       1,  // use default DB
	//
	//})
	//
	//// 设置限流
	//server.Use(ratelimit.NewBuilder(cmd, time.Minute, 10).Build())

	// 配置cors，解决跨域问题  配置一个middleware 作用于所有的方法 （类似Spring中的拦截器）
	server.Use(cors.New(cors.Config{

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
	}))

	// 下面使用了两个middleware，进行登录校验
	//store := cookie.NewStore([]byte("secret"))

	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"), []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))

	//store, err := cookie.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"), []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
	//if err != nil {
	//	panic(err)
	//}

	//server.Use(sessions.Sessions("mysession", store))

	//注册登录状态middleware
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup").Build())

	// 使用JWT做验证
	useJWT(server)

	// middleware相当于spring中的拦截器，起到一个请求前置处理的作用 aop

	return server
}

func useJWT(server *gin.Engine) {
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/loginJWT").
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())
}

// 初始化用户 无自动注入
func initUser(db *gorm.DB) *web.UserHandler {
	// dao
	ud := dao.NewUserDAO(db)

	// redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
	}

	userCache := cache.NewUserCache(rdb)

	// repository
	repo := repository.NewUserRepository(ud, userCache)

	// service
	svc := service.NewUserService(repo)

	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)

	// controller
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

// 初始化数据库
func initDB() *gorm.DB {
	// mysql连接地址 这里只是本地环境的连法
	db, err := gorm.Open(mysql.Open("root@tcp(localhost:3306)/webook"))

	// mysql连接地址，这里是k8s内部的连接法，pod与pod之间通过service的name和port进行连接
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-live-mysql:11309)/webook"))
	if err != nil {
		// panic表示该goroutine直接结束。只在main函数的初始化过程中，使用panic
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	// 初始化表结构
	err = dao.InitTable(db)

	if err != nil {
		panic(err)
	}

	return db
}
