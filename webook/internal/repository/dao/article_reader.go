package dao

import (
	"context"

	"gorm.io/gorm"
)

type ReaderDAO interface {
	// Upsert 插入或者更新
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

type ReaderGormDAO struct {
	db *gorm.DB
}

func NewReaderGormDAO(db *gorm.DB) *ReaderGormDAO {
	return &ReaderGormDAO{
		db: db,
	}
}

func (dao *ReaderGormDAO) Upsert(ctx context.Context, art Article) error {
	return nil
}

func (dao *ReaderGormDAO) UpsertV2(ctx context.Context, art PublishedArticle) error {
	return nil
}
