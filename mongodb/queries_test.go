package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type MongoDBTestSuite struct {
	suite.Suite
	col *mongo.Collection
}

// SetupSuite 在测试套件开始前执行,用于初始化测试环境
func (s *MongoDBTestSuite) SetupSuite() {
	// 获取测试对象
	t := s.T()

	// 创建一个10秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保在函数退出时取消上下文

	// 创建MongoDB命令监视器,用于调试和监控
	// 每当执行MongoDB命令时,会打印出具体的命令内容
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}

	// 配置MongoDB客户端连接选项
	// 设置连接URI和命令监视器
	// URI格式: mongodb://用户名:密码@主机:端口/
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/").
		SetMonitor(monitor)

	// 连接到MongoDB服务器
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	// 获取webook数据库中的articles集合
	col := client.Database("webook").
		Collection("articles")
	s.col = col // 将集合保存到测试套件中供后续测试使用

	// 插入测试数据
	// 插入两条Article文档用于后续查询测试
	manyRes, err := col.InsertMany(ctx, []any{Article{
		Id:       123,
		AuthorId: 11,
	}, Article{
		Id:       234,
		AuthorId: 12,
	}})

	// 确保插入操作没有错误
	assert.NoError(s.T(), err)
	// 记录实际插入的文档数量
	s.T().Log("插入数量", len(manyRes.InsertedIDs))
}

// TearDownSuite 在测试套件结束后执行,用于清理测试环境
func (s *MongoDBTestSuite) TearDownSuite() {
	// 创建一个1秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // 确保在函数退出时取消上下文

	// 删除集合中的所有文档
	// 使用空的bson.D{}作为过滤条件,表示匹配所有文档
	_, err := s.col.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)

	// 删除集合上的所有索引
	// 这确保下一次测试运行时有一个干净的环境
	_, err = s.col.Indexes().DropAll(ctx)
	assert.NoError(s.T(), err)
}

// TestOr 测试或查询
func (s *MongoDBTestSuite) TestOr() {
	// 创建一个1秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 定义查询条件
	// 使用bson.A创建一个包含两个条件的数组
	// 每个条件是一个bson.D,表示一个查询条件
	filter := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"id", 234}}}
	res, err := s.col.Find(ctx, bson.D{bson.E{"$or", filter}})
	assert.NoError(s.T(), err)

	// 定义一个用于存储查询结果的切片
	var arts []Article

	// 将查询结果解码到arts切片中
	err = res.All(ctx, &arts)
	assert.NoError(s.T(), err)
	s.T().Log("查询结果", arts)
}

// TestAnd 测试AND查询操作
func (s *MongoDBTestSuite) TestAnd() {
	// 创建一个1秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // 确保在函数退出时取消上下文

	// 构建AND查询条件
	// 使用bson.A创建一个包含两个条件的数组
	// 第一个条件:id=123
	// 第二个条件:author_id=11
	filter := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"author_id", 11}}}

	// 执行查询操作
	// 使用$and操作符组合多个查询条件
	res, err := s.col.Find(ctx, bson.D{bson.E{"$and", filter}})
	assert.NoError(s.T(), err)

	// 定义一个切片用于存储查询结果
	var arts []Article

	// 将查询结果解码到arts切片中
	err = res.All(ctx, &arts)
	assert.NoError(s.T(), err)

	// 打印查询结果
	s.T().Log("查询结果", arts)
}

// TestIn 测试IN查询操作
func (s *MongoDBTestSuite) TestIn() {
	// 创建一个1秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // 确保在函数退出时取消上下文

	// 构建IN查询条件
	// 使用bson.D创建一个查询条件,查询id在[123,234]范围内的文档
	filter := bson.D{bson.E{"id",
		bson.D{bson.E{"$in", []int{123, 234}}}}}

	// 定义投影,只返回id字段
	// 使用bson.M定义投影,1表示返回该字段
	proj := bson.M{"id": 1}

	// 执行查询操作
	// 使用SetProjection设置投影,只返回指定字段
	res, err := s.col.Find(ctx, filter,
		options.Find().SetProjection(proj))
	assert.NoError(s.T(), err)

	// 定义一个切片用于存储查询结果
	var arts []Article

	// 将查询结果解码到arts切片中
	err = res.All(ctx, &arts)
	assert.NoError(s.T(), err)

	// 打印查询结果
	s.T().Log("查询结果", arts)
}

// TestIndexes 测试创建MongoDB索引
// 该测试函数演示了如何在MongoDB集合上创建唯一索引
// 主要步骤:
// 1. 创建带超时的上下文
// 2. 构建索引模型,设置索引键和选项
// 3. 调用CreateOne创建单个索引
// 4. 验证索引创建结果
func (s *MongoDBTestSuite) TestIndexes() {
	// 创建1秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保函数退出时取消上下文
	defer cancel()

	// 创建索引
	// - Keys: 指定在id字段上创建升序(1)索引
	// - Options: 设置为唯一索引,并指定索引名称为"idx_id"
	ires, err := s.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{"id", 1}},
		Options: options.Index().SetUnique(true).SetName("idx_id"),
	})
	// 确保索引创建成功,无错误发生
	assert.NoError(s.T(), err)
	// 打印创建的索引名称
	s.T().Log("创建索引", ires)
}

func TestMongoDBQueries(t *testing.T) {
	suite.Run(t, &MongoDBTestSuite{})
}
