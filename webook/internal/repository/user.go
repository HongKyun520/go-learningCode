package repository

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository/dao"
	"context"
)

// 在repository层面，就没有业务上的概念了，有的只是对数据库的curd的概念
var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

// 构造函数，实例永远是从外侧传入
func NewUserRepository(d *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: d}
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
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(int642 int64) {

}
