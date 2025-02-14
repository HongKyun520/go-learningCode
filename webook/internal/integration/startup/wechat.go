package startup

import (
	"GoInAction/webook/internal/service/oauth2/wechat"
	"GoInAction/webook/pkg/logger"
)

func InitWechatService(l logger.Logger) wechat.Service {
	return wechat.NewWechatService("", "")
}
