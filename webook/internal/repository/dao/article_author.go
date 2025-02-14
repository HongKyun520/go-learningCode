package dao

import (
	"context"

	"gorm.io/gorm"
)

type AuthorDAO interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type AuthorGormDAO struct {
	db *gorm.DB
}

func (dao *AuthorGormDAO) Create(ctx context.Context, art Article) (int64, error) {
	return 0, nil
}

func (dao *AuthorGormDAO) Update(ctx context.Context, art Article) error {
	return nil
}

func NewAuthorGormDAO(db *gorm.DB) *AuthorGormDAO {
	return &AuthorGormDAO{
		db: db,
	}
}
