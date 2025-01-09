//go:build wireinject

package wire

import (
	"GoInAction/wire/repository"
	"GoInAction/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {

	// 使用wire.Build来构建依赖关系
	// repository.NewUserRepository - 创建UserRepository实例
	// dao.NewUserDAO - 创建UserDAO实例
	// InitDB - 初始化数据库连接
	// 这些方法会按照依赖顺序被wire自动调用
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	// 返回空对象,wire会在编译时替换为真实实例
	return new(repository.UserRepository)

}
