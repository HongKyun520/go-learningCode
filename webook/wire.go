//go:build wireinject

package main

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/web"
	"GoInAction/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		ioc.InitSmsService,
		service.NewUserService,
		service.NewCodeService,
		web.NewUserHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return gin.Default()
}
