package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

type Address struct {
	Id     int64
	UserId int64
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数 / 存纳秒数
	milli := time.Now().UnixMilli()
	u.Utime = milli
	u.Ctime = milli
	// 使用gorm插入一条数据，并将链路保持下去

	// 获取mysql数据库的错误码
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == 1062 {
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err1 := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	err1 = dao.db.WithContext(ctx).First(&u, "email = ?", email).Error

	return u, err1
}

// DAO层直接定义po对象，对应数据库表
// 有些字段在数据库是JSON格式存储的，那么在domain里面就会被转为结构体
// DAO层的字段需要加上gorm的定义
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全局唯一索引
	Email    string `gorm:"unique"`
	Password string

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
