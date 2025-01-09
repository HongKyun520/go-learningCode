package ioc

import (
	"GoInAction/webook/internal/service/sms"
	"GoInAction/webook/internal/service/sms/memory"
	// "GoInAction/webook/internal/service/sms/tencent"
)

func InitSmsService() sms.Service {
	return memory.NewService()
}

// func InitSmsServiceTencent() sms.Service {

// 	return tencent.NewService()
// }
