package dao

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

// TestS3 测试腾讯云对象存储COS的基本操作
// 主要测试以下功能:
// 1. 从环境变量获取COS的访问凭证
// 2. 初始化AWS S3会话配置
// 3. 上传对象到COS存储桶
// 4. 从COS存储桶下载对象
func TestS3(t *testing.T) {
	// 腾讯云中对标 s3 和 OSS 的产品叫做 COS
	// 从环境变量获取COS应用ID
	cosId, ok := os.LookupEnv("COS_APP_ID")
	if !ok {
		panic("没有找到环境变量 COS_APP_ID ")
	}
	// 从环境变量获取COS应用密钥
	cosKey, ok := os.LookupEnv("COS_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 COS_APP_SECRET")
	}

	// 初始化AWS S3会话配置
	sess, err := session.NewSession(&aws.Config{
		// 设置访问凭证
		Credentials: credentials.NewStaticCredentials(cosId, cosKey, ""),
		// 设置COS区域
		Region: ekit.ToPtr[string]("ap-nanjing"),
		// 设置COS访问端点
		Endpoint: ekit.ToPtr[string]("https://cos.ap-nanjing.myqcloud.com"),
		// 强制使用 /bucket/key 的形态
		S3ForcePathStyle: ekit.ToPtr[bool](true),
	})
	assert.NoError(t, err)

	// 创建S3客户端
	client := s3.New(sess)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 上传对象到COS
	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		// 指定存储桶名称
		Bucket: ekit.ToPtr[string]("webook-1314583317"),
		// 指定对象键
		Key: ekit.ToPtr[string]("12634"),
		// 设置对象内容
		Body: bytes.NewReader([]byte("测试内容 abc")),
		// 设置内容类型
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	assert.NoError(t, err)

	// 从COS下载对象
	res, err := client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		// 指定存储桶名称
		Bucket: ekit.ToPtr[string]("webook-1314583317"),
		// 指定对象键
		Key: ekit.ToPtr[string]("测试文件"),
	})
	assert.NoError(t, err)

	// 读取对象内容
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	// 打印对象内容
	t.Log(string(data))
}
