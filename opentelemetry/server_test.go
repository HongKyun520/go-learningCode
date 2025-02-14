package opentelemetry

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// TestServer 测试OpenTelemetry的基本功能
// 包括:
// 1. 创建和配置资源(Resource)
// 2. 设置传播器(Propagator)用于在服务间传递追踪上下文
// 3. 配置追踪提供者(TracerProvider)
// 4. 启动HTTP服务器并创建示例追踪
func TestServer(t *testing.T) {
	// 创建一个新的资源,包含服务名称和版本信息
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	// 创建传播器并设置为全局传播器
	prop := newPropagator()
	// 在客户端和服务端之间传递 tracing 的相关信息
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	// 创建Gin服务器
	server := gin.Default()
	server.GET("/test", func(ginCtx *gin.Context) {
		// 创建一个tracer,名字需要唯一
		tracer := otel.Tracer("gitee.com/geekbang/basic-go/opentelemetry")
		var ctx context.Context = ginCtx
		// 创建顶层span
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()
		time.Sleep(time.Second)
		// 添加事件
		span.AddEvent("发生了某事")
		// 创建子span
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()
		// 设置span属性
		subSpan.SetAttributes(attribute.String("attr1", "value1"))
		time.Sleep(time.Millisecond * 300)
		ginCtx.String(http.StatusOK, "测试 span")
	})
	server.Run(":8082")
}

// newResource 创建一个新的Resource
// 参数:
//   - serviceName: 服务名称
//   - serviceVersion: 服务版本
//
// 返回:
//   - *resource.Resource: 创建的资源
//   - error: 错误信息
func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

// newTraceProvider 创建一个新的TracerProvider
// 参数:
//   - res: 资源信息
//
// 返回:
//   - *trace.TracerProvider: 创建的追踪提供者
//   - error: 错误信息
func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	// 创建Zipkin导出器
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	// 创建追踪提供者
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// 默认为5秒,这里设置为1秒用于演示
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

// newPropagator 创建一个新的传播器
// 返回:
//   - propagation.TextMapPropagator: 创建的传播器
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
