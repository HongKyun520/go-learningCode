package dao

import (
	"GoInAction/webook/internal/domain"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, artId int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
}

type ArticleGormDAO struct {
	db *gorm.DB
}

func NewArticleGORMDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGormDAO{
		db: db,
	}
}

// SyncStatus 同步文章状态到文章表和已发布文章表
// 使用事务确保数据一致性,如果任何操作失败会回滚整个事务
//
// 参数:
//   - ctx: 上下文,用于控制超时和取消
//   - uid: 作者ID,用于验证文章归属
//   - artId: 要更新状态的文章ID
//   - status: 新的文章状态
//
// 返回:
//   - error: 错误信息。如果同步过程中出现任何错误,将回滚事务
//
// 处理流程:
//  1. 开启数据库事务
//  2. 更新文章表中对应文章的状态,同时验证作者ID
//  3. 检查是否成功更新了一行记录
//  4. 更新已发布文章表中的状态
//  5. 提交事务
func (dao *ArticleGormDAO) SyncStatus(ctx context.Context, uid int64, artId int64, status domain.ArticleStatus) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", artId, uid).Updates(map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		})
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("更新失败, 找不到该文章")
		}

		return tx.Model(&PublishedArticle{}).Where("id = ? ", artId).Update("status", status).Error
	})
}

func (dao *ArticleGormDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pubArt PublishedArticle
	err := dao.db.WithContext(ctx).Where("id = ?", id).Where("status = ?", domain.ArticleStatusPublished).First(&pubArt).Error
	return pubArt, err
}

func (dao *ArticleGormDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	return art, err
}

func (dao *ArticleGormDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&Article{}).Where("id = ? AND author_id", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   now,
		"status":  art.Status,
	})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 0 {
		return errors.New("更新失败")
	}

	return nil

}

func (dao *ArticleGormDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// Sync 同步文章数据到文章表和已发布文章表
// 使用事务确保数据一致性,如果任何操作失败会回滚整个事务
//
// 参数:
//   - ctx: 上下文,用于控制超时和取消
//   - art: 要同步的文章对象
//
// 返回:
//   - int64: 文章ID。如果是更新操作返回原ID,如果是创建操作返回新ID
//   - error: 错误信息。如果同步过程中出现任何错误,将回滚事务
//
// 处理流程:
//  1. 开启数据库事务
//  2. 根据文章是否存在ID判断是更新还是新建操作
//  3. 同步数据到文章表
//  4. 同步数据到已发布文章表,使用 upsert 操作
//  5. 提交事务
func (dao *ArticleGormDAO) Sync(ctx context.Context, art Article) (int64, error) {
	tx := dao.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer tx.Rollback()

	var (
		id  = art.Id
		err error
	)

	dao1 := NewArticleGORMDAO(tx)
	if id > 0 {
		err = dao1.UpdateById(ctx, art)
	} else {
		id, err = dao1.Insert(ctx, art)
	}

	if err != nil {
		return 0, err
	}

	art.Id = id
	now := time.Now().UnixMilli()
	pubArt := PublishedArticle(art)
	pubArt.Ctime = now
	pubArt.Utime = now
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"title": pubArt.Title, "content": pubArt.Content, "utime": now, "status": pubArt.Status}),
	}).Create(&pubArt).Error

	if err != nil {
		return 0, err
	}

	return id, tx.Commit().Error
}

func (dao *ArticleGormDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var arts []Article
	err := dao.db.WithContext(ctx).Where("author_id = ?", uid).Offset(offset).Limit(limit).Find(&arts).Order("utime DESC").Error
	return arts, err
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}

type PublishedArticle Article
