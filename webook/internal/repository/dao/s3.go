package dao

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleS3DAO struct {
	ArticleGormDAO
	oss *s3.S3
}

// SyncStatus 同步文章状态并处理对应的存储对象
//
// 主要功能:
// 1. 在数据库中同步更新文章状态
// 2. 如果文章状态变为私有,则删除对应的OSS存储对象
//
// 参数:
//   - ctx: 上下文,用于控制超时和取消
//   - uid: 用户ID,用于验证文章归属
//   - id: 文章ID
//   - status: 新的文章状态
//
// 处理流程:
// 1. 开启数据库事务
// 2. 更新文章表中对应文章的状态,同时验证作者ID
// 3. 更新已发布文章表中的状态
// 4. 如果状态为私有(3),则删除OSS中存储的文章内容
//
// 返回:
//   - error: 如果发生任何错误则返回对应错误信息
func (a *ArticleS3DAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", uid, id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("ID 不对或者创作者不对")
		}
		return tx.Model(&PublishedArticleV2{}).
			Where("id = ?", uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
	if err != nil {
		return err
	}
	const statusPrivate = 3
	if status == statusPrivate {
		_, err = a.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: ekit.ToPtr[string]("webook-1314583317"),
			Key:    ekit.ToPtr[string](strconv.FormatInt(id, 10)),
		})
	}
	return err
}

func (a *ArticleS3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		dao := NewArticleGORMDAO(tx)
		if id > 0 {
			err = dao.UpdateById(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		now := time.Now().UnixMilli()
		pubArt := PublishedArticleV2{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Ctime:    now,
			Utime:    now,
			Status:   art.Status,
		}
		pubArt.Ctime = now
		pubArt.Utime = now
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  pubArt.Title,
				"utime":  now,
				"status": pubArt.Status,
			}),
		}).Create(&pubArt).Error
		return err
	})
	if err != nil {
		return 0, err
	}
	_, err = a.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("webook-1314583317"),
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return id, err
}

func NewArticleS3DAO(db *gorm.DB, oss *s3.S3) *ArticleS3DAO {
	return &ArticleS3DAO{ArticleGormDAO: ArticleGormDAO{db: db}, oss: oss}
}

type PublishedArticleV2 struct {
	Id    int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}
