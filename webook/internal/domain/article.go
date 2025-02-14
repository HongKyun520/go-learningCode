package domain

import "time"

const (
	// ArticleStatusUnknown 这是一个未知状态
	ArticleStatusUnknown = iota
	// ArticleStatusUnpublished 未发表
	ArticleStatusUnpublished
	// ArticleStatusPublished 已发表
	ArticleStatusPublished
	// ArticleStatusPrivate 仅自己可见
	ArticleStatusPrivate
)

type ArticleStatus uint8

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

// 生成摘要，从内容中抽取前100个字符
func (a Article) Abstract() string {
	str := []rune(a.Content)
	if len(str) < 128 {
		return string(str)
	}
	return string(str[:128]) + "..."
}

func (a Article) ToUint8() uint8 {
	return uint8(a.Status)
}

type Author struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
