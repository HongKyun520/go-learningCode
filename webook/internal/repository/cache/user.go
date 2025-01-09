package cache

import (
	"GoInAction/webook/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheNotFound = errors.New("缓存不存在")
var ErrKeyNotFound = redis.Nil

// 业务层抽象，引入一个专门的UserCache解决

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	// 面向接口设计
	// 传单机 redis可以
	// 传cluster 的 redis也可以
	client redis.Cmdable

	// 过期时间
	expiration time.Duration
}

// A用到B，B一定是接口  => 保证面向接口
// A用到了B，B一定是A的字段  => 规避包变量、包方法，这两者都非常缺乏扩展性
// A用到了B，A绝对不初始化B，而是外面注入  => 保持依赖注入和依赖反转（DI 和 IOC）
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// get方法
// 只要err为nil，则缓存里有数据
// 缓存没有数据，则返回一个特定的err
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	u := domain.User{}
	key := cache.getKey(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return u, err
	}

	// 反序列化
	err = json.Unmarshal(val, &u)

	return u, err

}

func (cache *RedisUserCache) GetV1(ctx context.Context, id int64) (domain.User, error) {
	u := domain.User{}
	key := cache.getKey1(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return u, err
	}
	// 反序列化
	err = json.Unmarshal(val, &u)

	return u, err

}

// set方法
func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {

	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	key := cache.getKey(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

// set方法
func (cache *RedisUserCache) SetV1(ctx context.Context, u domain.User) error {

	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	key := cache.getKey(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) getKey(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *RedisUserCache) getKey1(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

// TODO: 本地缓存实现
