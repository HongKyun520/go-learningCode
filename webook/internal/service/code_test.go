package service

import (
	"GoInAction/webook/internal/repository"
	repomocks "GoInAction/webook/internal/repository/mocks"
	"GoInAction/webook/internal/service/sms"
	smsmocks "GoInAction/webook/internal/service/sms/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCodeService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service)
		biz     string
		phone   string
		wantErr error
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				smsSvc := smsmocks.NewMockService(ctrl)

				// 设置验证码存储期望
				repo.EXPECT().Store(gomock.Any(), "login", "13800138000", gomock.Any()).Return(nil)
				// 设置短信发送期望
				smsSvc.EXPECT().Send(gomock.Any(), "13800138000", gomock.Any()).Return(nil)

				return repo, smsSvc
			},
			biz:     "login",
			phone:   "13800138000",
			wantErr: nil,
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				smsSvc := smsmocks.NewMockService(ctrl)

				// 模拟存储时发现发送太频繁
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(repository.ErrCodeSendTooMany)

				return repo, smsSvc
			},
			biz:     "login",
			phone:   "13800138000",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误-存储失败",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				smsSvc := smsmocks.NewMockService(ctrl)

				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("redis error"))

				return repo, smsSvc
			},
			biz:     "login",
			phone:   "13800138000",
			wantErr: errors.New("redis error"),
		},
		{
			name: "系统错误-发送失败",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				smsSvc := smsmocks.NewMockService(ctrl)

				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				smsSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("sms error"))

				return repo, smsSvc
			},
			biz:     "login",
			phone:   "13800138000",
			wantErr: errors.New("sms error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo, smsSvc := tc.mock(ctrl)
			svc := NewCodeService(repo, smsSvc)

			err := svc.Send(context.Background(), tc.biz, tc.phone)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCodeService_Verify(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) repository.CodeRepository
		biz     string
		phone   string
		code    string
		wantOk  bool
		wantErr error
	}{
		{
			name: "验证成功",
			mock: func(ctrl *gomock.Controller) repository.CodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Verify(gomock.Any(), "login", "13800138000", "123456").
					Return(true, nil)
				return repo
			},
			biz:     "login",
			phone:   "13800138000",
			code:    "123456",
			wantOk:  true,
			wantErr: nil,
		},
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) repository.CodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, nil)
				return repo
			},
			biz:     "login",
			phone:   "13800138000",
			code:    "000000",
			wantOk:  false,
			wantErr: nil,
		},
		{
			name: "验证次数太多",
			mock: func(ctrl *gomock.Controller) repository.CodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, repository.ErrCodeVerifyTooManyTimes)
				return repo
			},
			biz:     "login",
			phone:   "13800138000",
			code:    "123456",
			wantOk:  false,
			wantErr: ErrCodeVerifyTooManyTimes,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.CodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, errors.New("redis error"))
				return repo
			},
			biz:     "login",
			phone:   "13800138000",
			code:    "123456",
			wantOk:  false,
			wantErr: errors.New("redis error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewCodeService(repo, nil)

			ok, err := svc.Verify(context.Background(), tc.biz, tc.code, tc.phone)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}
