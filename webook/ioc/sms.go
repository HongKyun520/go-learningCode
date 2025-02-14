package ioc

import (
	"GoInAction/webook/internal/service/sms"
	"GoInAction/webook/internal/service/sms/memory"
	// "GoInAction/webook/internal/service/sms/tencent"
)

func InitSmsService() sms.Service {
	// return ratelimit.NewRateLimitSMSService(memory.NewService(), limiter.NewRedisSlideWindowLimiter(redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// }), time.Second, 100), "sms:limit")
	return memory.NewService()
}

// func InitSmsServiceTencent() sms.Service {

// 	return tencent.NewService()
// }
