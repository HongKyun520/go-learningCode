package repository

import "GoInAction/wire/repository/dao"

type UserRepository struct {
	userDao *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{
		userDao: userDAO,
	}
}
