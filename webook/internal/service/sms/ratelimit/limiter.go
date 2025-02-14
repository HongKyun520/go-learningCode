package ratelimit

import (
	"GoInAction/webook/internal/service/sms"
	"GoInAction/webook/pkg/limiter"
	"context"
	"errors"
	"log"
)

var errLimit = errors.New("send sms too many times")

// RateLimitSMSService 带限流功能的短信服务 装饰器实现
type RateLimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter, key string) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
		key:     key,
	}
}

// Send 发送短信,带限流功能
// ctx 上下文
// tplId 短信模板ID
// args 短信模板参数
// phone 接收短信的手机号列表
// 返回error,如果限流或发送失败则返回错误
func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, phone ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		log.Println("limit error:", err)
		return err
	}
	if limited {
		log.Println("send sms too many times")
		return errLimit
	}
	return r.svc.Send(ctx, tplId, args, phone...)
}
