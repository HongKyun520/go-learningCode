package main

import (
	"GoInAction/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// 应用启动入口

func main() {
	server := gin.Default()

	// middleware相当于spring中的拦截器，起到一个请求前置处理的作用
	server.Use(func(context *gin.Context) {
		println("这是第一个middleware")
	})

	server.Use(func(context *gin.Context) {
		println("这是第二个middleware")
	})

	// 配置cros  配置一个middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:8080"},
		AllowMethods: []string{"PUT", "PATCH", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		// 自定义路由策略
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	u := web.NewUserHandler()
	u.RegisterUsersRoutes(server)

	server.Run(":8080")
}
