package main

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/repository/dao"
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/web"
	"GoInAction/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

// main方法需要自己手写很多东西

// go项目需要在main方法中手动做很多事情
// 初始化数据源、甚至初始化普通类、手动注册接口路由、手动注册middleware(过滤器、拦截器)

// 应用启动入口
func main() {
	// 需要手动初始化userHandler、userService、userRepository、userDao、mysqlDB ... 比较坑
	// 不像spring框架可以直接手动注入

	// 初始化数据库
	//db := initDB()

	// 初始化用户
	//u := initUser(db)

	// 初始化服务器
	//server := initWebServer()

	// 将userHandler里面的路由注册进server
	//u.RegisterUsersRoutes(server)
	//server.Run(":8080")

	engine := gin.Default()
	engine.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world")
	})

	engine.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// middleware相当于spring中的拦截器，起到一个请求前置处理的作用
	//server.Use(func(context *gin.Context) {
	//	println("这是第一个middleware")
	//}, func(context *gin.Context) {
	//	println("这是第二个middleware")
	//})

	// 配置cors，解决跨域问题  配置一个middleware 作用于所有的方法 （类似Spring中的拦截器）
	server.Use(cors.New(cors.Config{

		// 允许的请求来源（一般填写域名）
		//AllowOrigins: []string{"https://localhost:8080"},
		// 允许的请求方法 不写默认都支持

		//AllowMethods: []string{"PUT", "PATCH", "POST"},
		// 允许的请求头
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//ExposeHeaders:    []string{"Content-Length"},
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
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	//注册登录状态middleware
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").Build())

	return server
}

// 初始化用户
func initUser(db *gorm.DB) *web.UserHandler {
	// dao
	ud := dao.NewUserDAO(db)

	// repository
	repo := repository.NewUserRepository(ud)

	// service
	svc := service.NewUserService(repo)

	// controller
	u := web.NewUserHandler(svc)
	return u
}

// 初始化数据库
func initDB() *gorm.DB {
	// 明文输入？
	db, err := gorm.Open(mysql.Open("root:Kyun1024!@tcp(localhost:3306)/webook"))
	if err != nil {
		// panic表示该goroutine直接结束。只在main函数的初始化过程中，使用panic
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
