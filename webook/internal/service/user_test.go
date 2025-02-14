package service

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository"
	repomocks "GoInAction/webook/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestCachedUserService_SignUp(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) repository.UserRepository
		user    domain.User
		wantErr error
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return repo
			},
			user: domain.User{
				Email:    "test@qq.com",
				Password: "hello#world123",
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repository.ErrUserDuplicateEmail)
				return repo
			},
			user: domain.User{
				Email:    "exists@qq.com",
				Password: "hello#world123",
			},
			wantErr: repository.ErrUserDuplicateEmail,
		},
		{
			name: "数据库错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
				return repo
			},
			user: domain.User{
				Email:    "test@qq.com",
				Password: "hello#world123",
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			err := svc.SignUp(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCachedUserService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@qq.com").Return(domain.User{
					Id:       123,
					Email:    "test@qq.com",
					Password: string(hashedPwd),
				}, nil)
				return repo
			},
			email:    "test@qq.com",
			password: "hello#world123",
			wantUser: domain.User{
				Id:    123,
				Email: "test@qq.com",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "notexist@qq.com",
			password: "hello#world123",
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Password: string(hashedPwd),
				}, nil)
				return repo
			},
			email:    "test@qq.com",
			password: "wrongpassword",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser.Id, user.Id)
				assert.Equal(t, tc.wantUser.Email, user.Email)
			}
		})
	}
}

func TestCachedUserService_Profile(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查询成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), int64(123)).Return(domain.User{
					Id:    123,
					Email: "test@qq.com",
				}, nil)
				return repo
			},
			id: 123,
			wantUser: domain.User{
				Id:    123,
				Email: "test@qq.com",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			id:       0,
			wantUser: domain.User{},
			wantErr:  repository.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.Profile(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser, user)
			}
		})
	}
}

func TestCachedUserService_FindOrCreate(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		phone    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "已存在用户",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), "13800138000").Return(domain.User{
					Id:    123,
					Phone: "13800138000",
				}, nil)
				return repo
			},
			phone: "13800138000",
			wantUser: domain.User{
				Id:    123,
				Phone: "13800138000",
			},
			wantErr: nil,
		},
		{
			name: "新建用户",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 第一次查询不存在
				repo.EXPECT().FindByPhone(gomock.Any(), "13800138000").Return(domain.User{}, repository.ErrUserNotFound)
				// 创建用户
				repo.EXPECT().Create(gomock.Any(), domain.User{
					Phone: "13800138000",
				}).Return(nil)
				// 再次查询返回创建的用户
				repo.EXPECT().FindByPhone(gomock.Any(), "13800138000").Return(domain.User{
					Id:    123,
					Phone: "13800138000",
				}, nil)
				return repo
			},
			phone: "13800138000",
			wantUser: domain.User{
				Id:    123,
				Phone: "13800138000",
			},
			wantErr: nil,
		},
		{
			name: "创建用户失败",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
				return repo
			},
			phone:    "13800138000",
			wantUser: domain.User{},
			wantErr:  errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo)

			user, err := svc.FindOrCreate(context.Background(), tc.phone)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser, user)
			}
		})
	}
}
