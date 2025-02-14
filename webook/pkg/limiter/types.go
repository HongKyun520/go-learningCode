package limiter

import "context"

type Limiter interface {
	// 限流
	Limit(ctx context.Context, key string) (bool, error)
}
