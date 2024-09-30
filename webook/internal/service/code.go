package service

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "123456"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (c *CodeService) Send(ctx context.Context,
	biz string,
	phone string) error {

	// 生成验证码
	code := c.generateCode()

	// set redis
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	// 发送短信
	err = c.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {

		// 这里怎么办?
	}

	return err
}

func (c *CodeService) Verify(ctx context.Context,
	biz string,
	inputCode string,
	phone string) (bool, error) {

	return c.repo.Verify(ctx, biz, phone, inputCode)
}

func (c *CodeService) VerifyV1(ctx context.Context,
	biz string,
	inputCode string,
	phone string) error {

	return nil
}

func (c *CodeService) generateCode() string {
	num := rand.Intn(1000000)

	return fmt.Sprintf("%6d", num)
}
