package service

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository"
	"GoInAction/webook/pkg/logger"
	"context"
	"errors"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, id int64, uid int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository

	// V1
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	logger     logger.Logger
}

func (s *articleService) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	return s.repo.GetPubById(ctx, id)
}

func (s *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return s.repo.GetById(ctx, id)
}

func (s *articleService) Withdraw(ctx context.Context, id int64, uid int64) error {
	art := domain.Article{
		Id:     id,
		Status: domain.ArticleStatusPrivate,
	}
	return s.repo.SyncStatus(ctx, uid, art.Id, domain.ArticleStatusPrivate)
}

func NewArticleServiceV1(repo repository.ArticleRepository,
	readerRepo repository.ArticleReaderRepository,
	authorRepo repository.ArticleAuthorRepository) ArticleService {
	return &articleService{
		repo:       repo,
		readerRepo: readerRepo,
		authorRepo: authorRepo,
	}
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (s *articleService) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return s.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (s *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return s.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	// 想到这里要先操作制作库
	// 这里操作线上库
	art.Status = domain.ArticleStatusPublished
	var (
		id  = art.Id
		err error
	)

	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		// 我可能线上库已经有数据了
		// 也可能没有
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			// 多接入一些 tracing 的工具
			a.logger.Error("保存到制作库成功但是到线上库失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		} else {
			return id, nil
		}
	}
	a.logger.Error("保存到制作库成功但是到线上库失败，重试耗尽",
		logger.Int64("aid", art.Id),
		logger.Error(err))
	return id, errors.New("保存到线上库失败，重试次数耗尽")
}

func (s *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := s.repo.Update(ctx, art)
		return art.Id, err
	}
	return s.repo.Create(ctx, art)
}
