package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// 初始化简表过程
func InitTables(db *gorm.DB) error {

	// 根据model自动建表,go不存在根据注解扫描的涉及
	return db.AutoMigrate(&User{}, &Article{})

}

// InitCollections 初始化MongoDB集合和索引
//
// 参数:
//   - mdb: MongoDB数据库连接实例
//
// 返回:
//   - error: 如果创建索引过程中发生错误则返回错误信息
//
// 功能:
//  1. 为articles集合创建索引:
//     - 在id字段上创建唯一索引
//     - 在author_id字段上创建普通索引
//  2. 为published_articles集合创建索引:
//     - 在id字段上创建唯一索引
//     - 在author_id字段上创建普通索引
//
// 说明:
//   - 使用1秒超时的context来控制索引创建操作
//   - 如果任一索引创建失败则返回错误
func InitCollections(mdb *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	col := mdb.Collection("articles")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{"author_id", 1}},
		},
	})
	if err != nil {
		return err
	}
	liveCol := mdb.Collection("published_articles")
	_, err = liveCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{"author_id", 1}},
		},
	})
	return err
}
