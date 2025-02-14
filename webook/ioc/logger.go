package ioc

import (
	"GoInAction/webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitLogger 初始化日志组件
// 使用zap作为日志库,将其封装为自定义的Logger接口
//
// 实现步骤:
// 1. 创建开发环境的zap配置
// 2. 从viper中读取log配置并解析到zap配置中
// 3. 使用配置构建zap logger实例
// 4. 将zap logger封装为ZapLogger返回
//
// 参数:
//   - 无
//
// 返回值:
//   - *logger.ZapLogger: 封装了zap.Logger的自定义logger实例
//
// 如果配置解析或logger创建失败会panic
func InitLogger() logger.Logger {

	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(zapLogger)
}
