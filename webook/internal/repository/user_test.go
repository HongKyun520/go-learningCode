package repository

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository/cache"
	cachemocks "GoInAction/webook/internal/repository/cache/mocks"
	"GoInAction/webook/internal/repository/dao"
	daomocks "GoInAction/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_Create(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		user    domain.User
		wantErr error
	}{
		{
			name: "创建成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().Insert(gomock.Any(), dao.User{
					Email: sql.NullString{
						String: "test@qq.com",
						Valid:  true,
					},
					Password: "123456",
				}).Return(nil)

				return userDAO, userCache
			},
			user: domain.User{
				Email:    "test@qq.com",
				Password: "123456",
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().Insert(gomock.Any(), gomock.Any()).
					Return(dao.ErrUserDuplicateEmail)

				return userDAO, userCache
			},
			user: domain.User{
				Email:    "exists@qq.com",
				Password: "123456",
			},
			wantErr: ErrUserDuplicateEmail,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDAO, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDAO, userCache)

			err := repo.Create(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCachedUserRepository_FindByEmail(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		email    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().FindByEmail(gomock.Any(), "test@qq.com").
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "test@qq.com",
							Valid:  true,
						},
						Password: "123456",
					}, nil)

				return userDAO, userCache
			},
			email: "test@qq.com",
			wantUser: domain.User{
				Id:       123,
				Email:    "test@qq.com",
				Password: "123456",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(dao.User{}, sql.ErrNoRows)

				return userDAO, userCache
			},
			email:    "notexist@qq.com",
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDAO, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDAO, userCache)

			user, err := repo.FindByEmail(context.Background(), tc.email)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser.Id, user.Id)
				assert.Equal(t, tc.wantUser.Email, user.Email)
				assert.Equal(t, tc.wantUser.Password, user.Password)
			}
		})
	}
}

func TestCachedUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "test@qq.com",
						Password: "123456",
					}, nil)

				return userDAO, userCache
			},
			id: 123,
			wantUser: domain.User{
				Id:       123,
				Email:    "test@qq.com",
				Password: "123456",
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotFound)

				userDAO.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "test@qq.com",
							Valid:  true,
						},
						Password: "123456",
					}, nil)

				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "test@qq.com",
					Password: "123456",
				}).Return(nil)

				return userDAO, userCache
			},
			id: 123,
			wantUser: domain.User{
				Id:       123,
				Email:    "test@qq.com",
				Password: "123456",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(domain.User{}, cache.ErrKeyNotFound)

				userDAO.EXPECT().FindById(gomock.Any(), gomock.Any()).
					Return(dao.User{}, sql.ErrNoRows)

				return userDAO, userCache
			},
			id:       0,
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDAO, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDAO, userCache)

			user, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser, user)
			}
		})
	}
}

func TestCachedUserRepository_FindByPhone(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		phone    string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().FindByPhone(gomock.Any(), "13800138000").
					Return(dao.User{
						Id: 123,
						Phone: sql.NullString{
							String: "13800138000",
							Valid:  true,
						},
					}, nil)

				return userDAO, userCache
			},
			phone: "13800138000",
			wantUser: domain.User{
				Id:    123,
				Phone: "13800138000",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDAO := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userDAO.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(dao.User{}, sql.ErrNoRows)

				return userDAO, userCache
			},
			phone:    "13800138000",
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userDAO, userCache := tc.mock(ctrl)
			repo := NewUserRepository(userDAO, userCache)

			user, err := repo.FindByPhone(context.Background(), tc.phone)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantUser, user)
			}
		})
	}
}
