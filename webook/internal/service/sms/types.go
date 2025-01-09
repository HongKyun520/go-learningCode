package sms

import "context"

// sms 服务接口
type Service interface {

	// 发送短信
	Send(ctx context.Context, tql string, args []string, phone ...string) error
}
