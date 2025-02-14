package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = errors.New("该用户不存在")
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type gormUserDAO struct {
	db *gorm.DB
}

type Address struct {
	Id     int64
	UserId int64
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &gormUserDAO{db: db}
}

func (dao *gormUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数 / 存纳秒数
	milli := time.Now().UnixMilli()
	u.Utime = milli
	u.Ctime = milli
	u.Profile = "{}"
	// 使用gorm插入一条数据，并将链路保持下去

	// 获取mysql数据库的错误码 唯一键冲突
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == 1062 {
			return ErrUserDuplicateEmail
		}
	}

	return err
}

// 根据phone来查
func (dao *gormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err1 := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	err1 = dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return u, ErrUserNotFound
	}

	return u, err1
}

func (dao *gormUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err1 := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&u).Error
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return u, ErrUserNotFound
	}

	return u, err1
}

func (dao *gormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err1 := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	err1 = dao.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return u, ErrUserNotFound
	}

	return u, err1
}

func (dao *gormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err1 := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	err1 = dao.db.WithContext(ctx).First(&u, "`id` = ?", id).Error
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return u, ErrUserNotFound
	}

	return u, err1
}

// DAO层直接定义po对象，对应数据库表
// 有些字段在数据库是JSON格式存储的，那么在domain里面就会被转为结构体
// DAO层的字段需要加上gorm的定义
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全局唯一索引
	Email    sql.NullString `gorm:"unique"`
	Password string
	NickName string
	// 唯一索引允许有多个空值
	// 但是不能有多个 “”
	Phone    sql.NullString `gorm:"unique"`
	Birthday string
	Profile  string `gorm:"type:json"`

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64

	// 根据查询的实际情况，建立索引

	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
}
