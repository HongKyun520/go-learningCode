package repository

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/repository/cache"
	"GoInAction/webook/internal/repository/dao"
	"context"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, artId int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	dao      dao.ArticleDAO
	reader   dao.ReaderDAO
	author   dao.AuthorDAO
	db       *gorm.DB
	cache    cache.ArticleCache
	userRepo UserRepository
}

func NewCacheArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository) ArticleRepository {
	return &CacheArticleRepository{
		dao:      dao,
		cache:    cache,
		userRepo: userRepo,
	}
}

func NewCacheArticleRepositoryV2(dao dao.ArticleDAO, reader dao.ReaderDAO, author dao.AuthorDAO, cache cache.ArticleCache, userRepo UserRepository) ArticleRepository {
	return &CacheArticleRepository{
		dao:      dao,
		reader:   reader,
		author:   author,
		cache:    cache,
		userRepo: userRepo,
	}
}

func (c *CacheArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {

	res, err := c.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 我现在要去查询对应的 User 信息，拿到创作者信息
	res = c.toDomain(dao.Article(art))
	author, err := c.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
		// 要额外记录日志，因为你吞掉了错误信息
		//return res, nil
	}
	res.Author.Name = author.NickName
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := c.cache.SetPub(ctx, res)
		if er != nil {
			// 记录日志
		}
	}()
	return res, nil
}

func (c *CacheArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {

	// 首先查询缓存
	res, err := c.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}

	art, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	res = c.toDomain(art)
	// 缓存回写
	go func() {
		err := c.cache.Set(ctx, res)
		if err != nil {
			// 记录日志
		}
	}()

	return res, nil
}

func (c *CacheArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 首先第一步，判定要不要查询缓存
	if offset == 0 && limit == 100 {
		arts, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return arts, err
		} else {
			// 记录日
		}
	}

	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})

	// 缓存回写
	go func() {
		// 使用一个新的context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if offset == 0 && limit == 100 {
			err := c.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
			}
		}
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		c.preCache(ctx, res)
	}()

	return res, nil
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, uid int64, artId int64, status domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, uid, artId, status)
	if err == nil {
		err = c.cache.DelFirstPage(ctx, uid)
		if err != nil {
			// 记录日志
		}
	}
	return err
}

func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err == nil {
		er := c.cache.DelFirstPage(ctx, art.Author.Id)
		if er != nil {
			// 也要记录日志
		}
	}
	// 在这里尝试，设置缓存当前页面的缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 你可以灵活设置过期时间
		user, er := c.userRepo.FindById(ctx, art.Author.Id)
		if er != nil {
			// 要记录日志
			return
		}
		art.Author = domain.Author{
			Id:   user.Id,
			Name: user.NickName,
		}
		er = c.cache.SetPub(ctx, art)
		if er != nil {
			// 记录日志
		}
	}()
	return id, nil
}

func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = c.author.Update(ctx, artn)
	} else {
		id, err = c.author.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = c.reader.Upsert(ctx, artn)
	return id, err
}

// SyncV2 同步文章数据到作者表和读者表，使用事务确保数据一致性
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - art: 要同步的文章领域对象
//
// 返回:
//   - int64: 文章ID。如果是更新操作返回原ID，如果是创建操作返回新ID
//   - error: 错误信息。如果同步过程中出现任何错误，将回滚事务
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	// 开启数据库事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// 确保在函数返回时回滚未提交的事务
	defer tx.Rollback()

	// 创建使用事务的 DAO 实例
	authorDAO := dao.NewAuthorGormDAO(tx)
	readerDAO := dao.NewReaderGormDAO(tx)

	// 将领域对象转换为持久化对象
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)

	// 根据是否存在 ID 决定是更新还是创建操作
	if id > 0 {
		// 更新已存在的文章
		err = authorDAO.Update(ctx, artn)
	} else {
		// 创建新文章
		id, err = authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}

	// 设置文章 ID 并同步到读者表
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(artn))

	// 提交事务
	tx.Commit()
	return id, err
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	err := c.dao.UpdateById(ctx, c.toEntity(art))
	if err == nil {
		err = c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			// 记录日志
		}
	}

	return err
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(art))
	if err == nil {
		err = c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			// 记录日志
		}
	}
	return id, err

}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

func (c *CacheArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			// 这里有一个错误
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}

// 预缓存
func (c *CacheArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 容量限制，太大的文章，不缓存
	const contentSizeThreshold = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		if err := c.cache.Set(ctx, arts[0]); err != nil {
			// 记录日志
		}

	}

}
