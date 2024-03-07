package domain

import "time"

// User是领域对象，是 DDD 中的 entity
// BO
type User struct {
	Id       int64
	Email    string `json:"email"`
	Password string `json:"password"`
	CTime    time.Time
}
