package tencent

import (
	"context"
	"fmt"
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

	return &Service{
		appId:     ekit.ToPtr[string](appId),
		signature: ekit.ToPtr[string](signName),
		client:    client,
	}

}

func (s Service) Send(ctx context.Context, tql string, args []string, phone ...string) error {

	req := sms.NewSendSmsRequest()

	req.SmsSdkAppId = s.appId
	req.SignName = s.signature

	req.TemplateId = ekit.ToPtr[string](tql)

	req.PhoneNumberSet = slice.Map[string, *string](phone, func(idx int, src string) *string {
		return &src
	})

	req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})

	response, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("send sms failed, code: %s, message: %s", *status.Code, *status.Message)
		}
	}

	return nil
}
