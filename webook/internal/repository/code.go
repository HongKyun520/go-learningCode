package repository

import (
	"GoInAction/webook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type cachedCodeRepository struct {
	// 引入cache层，而不是直接引入redis
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &cachedCodeRepository{
		cache: cache,
	}
}

func (c *cachedCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *cachedCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}
