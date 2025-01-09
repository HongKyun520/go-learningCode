package service

import (
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/internal/service/sms"
	"context"
	"fmt"
	"log"
	"math/rand"
)

const codeTplId = "123456"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context,
		biz string,
		phone string) error
	Verify(ctx context.Context,
		biz string,
		inputCode string,
		phone string) (bool, error)
}

type cachedCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &cachedCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// 发送验证码
func (c *cachedCodeService) Send(ctx context.Context,
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
		// 打印错误日志
		log.Println("send sms failed, err: ", err)
	}

	return err
}

// 验证验证码
func (c *cachedCodeService) Verify(ctx context.Context,
	biz string,
	inputCode string,
	phone string) (bool, error) {

	return c.repo.Verify(ctx, biz, phone, inputCode)
}

func (c *cachedCodeService) VerifyV1(ctx context.Context,
	biz string,
	inputCode string,
	phone string) error {

	return nil
}

// 生成验证码
func (c *cachedCodeService) generateCode() string {
	num := rand.Intn(1000000)

	return fmt.Sprintf("%06d", num)
}
