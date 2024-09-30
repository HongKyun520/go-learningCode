package repository

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"context"
)

// 在repository层面，就没有业务上的概念了，有的只是对数据库的curd的概念
var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

// 构造函数，实例永远是从外侧传入
func NewUserRepository(d *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{dao: d, cache: cache}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	// 在这里需要操作缓存。
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

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

	user = domain.User{
		Id:    ue.Id,
		Email: ue.Email,
		//Password: ue.Password,
		Phone: ue.Phone,
	}

	//err = r.cache.Set(ctx, user)
	//if err != nil {
	//	// 这里设置缓存失败, 打印日志就行
	//}

	// 异步更新缓存
	go func() {
		err = r.cache.Set(ctx, user)
		if err != nil {
			// 这里设置缓存失败, 打印日志就行
		}
	}()

	return user, err
}
