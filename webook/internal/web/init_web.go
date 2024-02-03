package web

import "github.com/gin-gonic/gin"

func RegisterRoutes() *gin.Engine {
	server := gin.Default()

	// 注册用户相关的路由
	//registerUsersRoutes(server)

	return server
}
