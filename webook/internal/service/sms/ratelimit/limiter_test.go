package ratelimit

import (
	"GoInAction/webook/internal/service/sms"
	smsmocks "GoInAction/webook/internal/service/sms/mocks"
	"GoInAction/webook/pkg/limiter"
	limitermocks "GoInAction/webook/pkg/limiter/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		tplId   string
		args    []string
		phones  []string
		wantErr error
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				// 不限流
				l.EXPECT().Limit(gomock.Any(), "test-key").Return(false, nil)
				// 发送成功
				svc.EXPECT().Send(gomock.Any(), "test-tpl", []string{"123456"}, "13800138000").
					Return(nil)

				return svc, l
			},
			tplId:   "test-tpl",
			args:    []string{"123456"},
			phones:  []string{"13800138000"},
			wantErr: nil,
		},
		{
			name: "被限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				// 触发限流
				l.EXPECT().Limit(gomock.Any(), "test-key").Return(true, nil)

				return svc, l
			},
			tplId:   "test-tpl",
			args:    []string{"123456"},
			phones:  []string{"13800138000"},
			wantErr: errLimit,
		},
		{
			name: "限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				// 限流器返回错误
				l.EXPECT().Limit(gomock.Any(), "test-key").
					Return(false, errors.New("redis error"))

				return svc, l
			},
			tplId:   "test-tpl",
			args:    []string{"123456"},
			phones:  []string{"13800138000"},
			wantErr: errors.New("redis error"),
		},
		{
			name: "发送短信失败",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				// 不限流
				l.EXPECT().Limit(gomock.Any(), "test-key").Return(false, nil)
				// 发送失败
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("sms error"))

				return svc, l
			},
			tplId:   "test-tpl",
			args:    []string{"123456"},
			phones:  []string{"13800138000"},
			wantErr: errors.New("sms error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, l := tc.mock(ctrl)
			limiter := NewRateLimitSMSService(svc, l, "test-key")

			err := limiter.Send(context.Background(), tc.tplId, tc.args, tc.phones...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
