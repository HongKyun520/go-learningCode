package tencent

import (
	"context"
	"fmt"
	"log"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId *string

	signature *string

	client *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {

	// 返回一个新的Service实例
	// appId: 腾讯云短信服务的应用ID
	// signName: 短信签名
	// client: 腾讯云短信服务的客户端实例
	return &Service{
		appId:     ekit.ToPtr[string](appId),    // 将appId转换为字符串指针
		signature: ekit.ToPtr[string](signName), // 将短信签名转换为字符串指针
		client:    client,                       // 设置短信客户端
	}

}

// Send 发送短信
// ctx 上下文
// tql 短信模板ID
// args 短信模板参数
// phone 接收短信的手机号列表
func (s *Service) Send(ctx context.Context, tql string, args []string, phone ...string) error {

	// 创建发送短信请求
	// 创建发送短信请求
	req := sms.NewSendSmsRequest()

	// 设置SDK AppID
	req.SmsSdkAppId = s.appId
	// 设置短信签名
	req.SignName = s.signature

	// 设置短信模板ID
	req.TemplateId = ekit.ToPtr[string](tql)

	// 设置接收短信的手机号列表
	req.PhoneNumberSet = slice.Map[string, *string](phone, func(idx int, src string) *string {
		return &src
	})

	// 设置短信模板参数
	req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})

	// 发送短信
	response, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	// 检查发送状态
	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("send sms failed, code: %s, message: %s", *status.Code, *status.Message)
		}
		log.Printf("send sms success, code: %s, message: %s", *status.Code, *status.Message)
	}

	return nil
}
