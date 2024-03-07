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

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// 保持跟handler侧的方法命名
// 不清楚返回什么的时候
// 返回一个error就行了
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {

	// 对密码进行加密赋值
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(pwd)
	return svc.repo.Create(ctx, u)

}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 找用户
	u, err := svc.repo.FindByEmail(ctx, email)

	// 找不到用户
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
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
