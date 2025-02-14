package repository

import (
	"GoInAction/webook/internal/domain"
	"context"
)

type ArticleReaderRepository interface {
	// Save 有则更新，无则插入，也就是 insert or update 语义
	Save(ctx context.Context, art domain.Article) error
}
