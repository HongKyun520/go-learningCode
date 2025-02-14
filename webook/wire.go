//go:build wireinject

package main

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/web"
	ijwt "GoInAction/webook/internal/web/jwt"
	"GoInAction/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO,
// 	cache.NewInteractiveRedisCache,
// 	repository.NewCachedInteractiveRepository,
// 	service.NewInteractiveService,
// )

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		// DAO部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		// interactiveSvcSet,

		// Cache部分
		cache.NewUserCache, cache.NewCodeCache,
		cache.NewArticleRedisCache,

		// Repository部分
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewCacheArticleRepository,

		// Service部分
		ioc.InitSmsService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// handler部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return gin.Default()
}
