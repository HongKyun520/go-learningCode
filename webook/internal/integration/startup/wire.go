//go:build wireinject

package startup

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

// thirdPartySet 定义了第三方依赖的提供者集合
// 包括:
// - Redis 客户端初始化
// - 数据库连接初始化
// - 日志组件初始化
var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,
	InitLogger)

// userSvcProvider 定义了用户服务相关的提供者集合
// 包括:
// - 用户数据访问对象(DAO)
// - 用户缓存
// - 用户仓储
// - 用户服务
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

// articlSvcProvider 定义了文章服务相关的提供者集合
// 包括:
// - 文章仓储
// - 文章Redis缓存
// - 文章数据访问对象
// - 文章服务
var articlSvcProvider = wire.NewSet(
	repository.NewCacheArticleRepository,
	cache.NewArticleRedisCache,
	dao.NewArticleGORMDAO,
	service.NewArticleService)

// InitWebServer 初始化并构建完整的Web服务器
// 使用wire进行依赖注入,组装所有组件
// 返回配置完成的gin.Engine实例
func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		articlSvcProvider,
		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSmsService,
		service.NewCodeService,
		InitWechatService,

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

// InitArticleHandler 初始化文章处理器
// 参数:
//   - dao: 文章数据访问对象
//
// 返回:
//   - *web.ArticleHandler: 文章处理器实例
func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		// interactiveSvcSet,
		repository.NewCacheArticleRepository,
		cache.NewArticleRedisCache,
		service.NewArticleService,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}
