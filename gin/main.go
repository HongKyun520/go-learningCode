package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建server
	server := gin.Default()

	// 注册路由
	server.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	// 静态路由
	server.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post")
	})

	// 参数路由 类似Java的@PathVariable
	// localhost:8080/user/2
	server.GET("/user/:name", func(c *gin.Context) {
		// 取到参数
		param := c.Param("name")
		c.String(http.StatusOK, "hello, 这是参数路由"+param)
	})

	// 通配符路由
	server.GET("/views/*.html", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, 通配符路由")
	})

	// 获取查询参数
	// localhost:8080/user/?id=123
	server.GET("/user", func(c *gin.Context) {
		value := c.Query("id")
		c.String(http.StatusOK, "hello, 这是查询参数"+value)
	})

	// 显式指定端口
	server.Run(":8080") // 监听并在0.0.0.0:8080 上启动服务
}
