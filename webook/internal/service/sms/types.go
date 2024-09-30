package sms

import "context"

type Service interface {

	// 发送短信
	Send(ctx context.Context, tql string, args []string, phone ...string) error
}
