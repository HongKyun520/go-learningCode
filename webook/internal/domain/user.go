package domain

import "time"

// User是领域对象，是 DDD 中的 entity
// BO
type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string

	Phone    string `json:"phone"`
	NickName string `json:"nick_name"`

	CTime time.Time
}
