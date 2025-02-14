package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"GoInAction/webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, art []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, art domain.Article) error
}

type ArticleRedisCache struct {
	client redis.Cmdable
}

func NewArticleRedisCache(client redis.Cmdable) ArticleCache {
	return &ArticleRedisCache{
		client: client,
	}
}

func (a *ArticleRedisCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	val, err := a.client.Get(ctx, a.pubKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (a *ArticleRedisCache) SetPub(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, a.pubKey(art.Id), val, time.Minute*10).Err()
}

func (a *ArticleRedisCache) Set(ctx context.Context, art domain.Article) error {
	key := a.key(art.Id)
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, data, time.Minute*10).Err()
}

func (a *ArticleRedisCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	key := a.key(id)
	artStr, err := a.client.Get(ctx, key).Result()
	if err != nil {
		return domain.Article{}, err
	}

	var art domain.Article
	err = json.Unmarshal([]byte(artStr), &art)
	if err != nil {
		return domain.Article{}, err
	}
	return art, nil
}

func (a *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := a.firstKey(uid)
	art, err := a.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal([]byte(art), &arts)
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (a *ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64, art []domain.Article) error {
	// 仅缓存摘要
	for i := range art {
		art[i].Content = art[i].Abstract()
	}

	key := a.firstKey(uid)
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, data, time.Minute*10).Err()
}

func (a *ArticleRedisCache) DelFirstPage(ctx context.Context, uid int64) error {
	key := a.firstKey(uid)
	return a.client.Del(ctx, key).Err()
}

func (a *ArticleRedisCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:detail:%d", id)
}

func (a *ArticleRedisCache) key(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

func (a *ArticleRedisCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
