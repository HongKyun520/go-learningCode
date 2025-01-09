package service

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 异常定义
var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
var ErrUserNotFound = errors.New("用户邮箱不存在")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type cachedUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &cachedUserService{repo: repo}
}

// 保持跟handler侧的方法命名
// 不清楚返回什么的时候
// 返回一个error就行了
func (svc *cachedUserService) SignUp(ctx context.Context, u domain.User) error {

	// 对密码进行加密赋值
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(pwd)
	return svc.repo.Create(ctx, u)
}

func (svc *cachedUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 找用户
	u, err := svc.repo.FindByEmail(ctx, email)

	// 找不到用户
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, err
	}

	// 比较密码，已加密，原始加密
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, err
}

// todo edit

// todo profile
func (svc *cachedUserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *cachedUserService) FindOrCreate(ctx context.Context,
	phone string) (domain.User, error) {

	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	// 没有这个用户
	u = domain.User{
		Phone: phone,
	}

	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}

	// 存在主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
	// return u, nil
}
