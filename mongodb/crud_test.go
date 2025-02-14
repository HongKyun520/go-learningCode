package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDB(t *testing.T) {
	// 创建一个带超时的上下文,设置10秒超时时间
	// 使用context.WithTimeout确保操作不会无限期挂起
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保在函数退出时取消上下文,避免资源泄露

	// 创建MongoDB命令监视器,用于调试和监控
	// 每当执行MongoDB命令时,会打印出具体的命令内容
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}

	// 配置MongoDB客户端连接选项
	// 设置连接URI和命令监视器
	// URI格式: mongodb://用户名:密码@主机:端口/认证数据库
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/admin").
		SetMonitor(monitor)

	// 使用配置的选项连接到MongoDB服务器
	// 返回一个MongoDB客户端实例
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	// 使用defer确保在测试结束时关闭MongoDB连接
	// 这是一个良好的实践,可以避免连接泄露
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			t.Error(err)
		}
	}()

	// 获取指定数据库(webook)和集合(articles)的引用
	// 后续的所有操作都将在这个集合上执行
	col := client.Database("webook").
		Collection("articles")

	// 插入一条测试文档
	// 使用Article结构体定义文档结构
	// InsertOne用于插入单个文档
	insertRes, err := col.InsertOne(ctx, Article{
		Id:       1,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 123,
	})
	assert.NoError(t, err)
	oid := insertRes.InsertedID.(primitive.ObjectID)
	t.Log("插入ID", oid)

	// 构建查询过滤器
	// 使用bson.M创建一个简单的查询条件
	// 这里查询id=1的文档
	filter := bson.M{
		"id": 1,
	}

	// 查找单个文档
	// FindOne返回符合条件的第一个文档
	// 使用Decode将文档解码到Article结构体中
	findRes := col.FindOne(ctx, filter)
	if findRes.Err() == mongo.ErrNoDocuments {
		t.Log("没找到数据")
	} else {
		assert.NoError(t, findRes.Err())
		var art Article
		err = findRes.Decode(&art)
		assert.NoError(t, err)
		t.Log(art)
	}

	// 构建更新操作的过滤器和更新内容
	// 使用bson.D保证字段顺序
	// $set操作符用于更新指定字段
	updateFilter := bson.D{bson.E{"id", 1}}
	set := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title": "新的标题",
	}}}

	// 更新单个文档
	// UpdateOne只更新符合条件的第一个文档
	// ModifiedCount表示实际更新的文档数量
	updateOneRes, err := col.UpdateOne(ctx, updateFilter, set)
	assert.NoError(t, err)
	t.Log("更新文档数量", updateOneRes.ModifiedCount)

	// 更新多个文档
	// UpdateMany会更新所有符合条件的文档
	// 这里使用整个Article对象作为更新内容
	updateManyRes, err := col.UpdateMany(ctx, updateFilter,
		bson.D{bson.E{Key: "$set",
			Value: Article{Content: "新的内容"}}})
	assert.NoError(t, err)
	t.Log("更新文档数量", updateManyRes.ModifiedCount)

	// 删除文档
	// DeleteMany删除所有符合条件的文档
	// DeletedCount表示实际删除的文档数量
	deleteFilter := bson.D{bson.E{"id", 1}}
	delRes, err := col.DeleteMany(ctx, deleteFilter)
	assert.NoError(t, err)
	t.Log("删除文档数量", delRes.DeletedCount)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}
