package repository

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"context"
	"database/sql"
	"log"
	"time"
)

// 在repository层面，就没有业务上的概念了，有的只是对数据库的curd的概念
var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type cachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

// 构造函数，实例永远是从外侧传入
func NewUserRepository(d dao.UserDAO, cache cache.UserCache) UserRepository {
	return &cachedUserRepository{dao: d, cache: cache}
}

func (r *cachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))

	// 在这里需要操作缓存。
}

func (r *cachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(u), nil
}

func (r *cachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	//
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(u), nil
}

func (r *cachedUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(u), nil
}

func (r *cachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	// 查缓存
	user, err := r.cache.Get(ctx, id)

	// 有数据
	if err == nil {
		return user, err
	}

	// 没数据，去数据库加载
	// 做好兜底策略
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user = r.entityToDomain(ue)
	// 异步更新缓存
	go func() {
		err = r.cache.Set(ctx, user)
		if err != nil {
			// 这里设置缓存失败, 打印日志就行
			log.Println(err)
		}
	}()

	return user, err
}

func (r *cachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		Ctime: u.CTime.UnixMilli(),
	}
}

func (r *cachedUserRepository) entityToDomain(ud dao.User) domain.User {
	return domain.User{
		Id:       ud.Id,
		Email:    ud.Email.String,
		Password: ud.Password,
		Phone:    ud.Phone.String,
		WechatInfo: domain.WechatInfo{
			UnionId: ud.WechatUnionId.String,
			OpenId:  ud.WechatOpenId.String,
		},
		CTime: time.UnixMilli(ud.Ctime),
	}
}
