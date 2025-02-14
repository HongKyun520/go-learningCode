package auth

import (
	"GoInAction/webook/internal/service/sms"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthSMSService struct {
	svc sms.Service
	key []byte
}

// Send 发送短信
// 使用JWT token验证短信模板权限
// ctx 上下文
// tplToken 短信模板token,包含了模板ID的JWT token
// args 短信模板参数
// numbers 接收短信的手机号列表
// 返回error,如果发送失败或token验证失败则返回错误
func (a *AuthSMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return err
	}
	return a.svc.Send(ctx, claims.TplId, args, numbers...)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	TplId string
}
