package ioc

import "GoInAction/webook/internal/service/oauth2/wechat"

func InitWechatService() wechat.Service {

	appID := "wx7256bc69ab349c72"
	appSecret := "1234567890"
	return wechat.NewWechatService(appID, appSecret)
}
