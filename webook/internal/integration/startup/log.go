package startup

import "GoInAction/webook/pkg/logger"

func InitLogger() logger.Logger {
	return logger.NewNopLogger()
}
