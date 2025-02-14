package prometheus

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type Builder struct {
	NameSpace  string
	Subsystem  string
	Name       string
	InstanceId string
}

// BuildResponseTime 构建一个用于监控HTTP请求响应时间的Gin中间件
// 该中间件会记录请求的响应时间,并通过Prometheus指标进行暴露
// 返回的指标名称格式为: {namespace}_{subsystem}_{name}_resp_time
//
// 记录的标签(labels)包括:
// - method: HTTP请求方法(GET/POST等)
// - pattern: 请求路径模式
// - status: HTTP响应状态码
// - instance_id: 实例ID(通过ConstLabels设置)
//
// 统计的分位数包括:
// - 0.5(中位数)
// - 0.75
// - 0.9
// - 0.99
// - 0.999
//
// 使用Summary类型的指标,可以计算请求响应时间的分布情况
func (b *Builder) BuildResponseTime() gin.HandlerFunc {

	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		// 不能有除_以外的符号
		Namespace: b.NameSpace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_resp_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,   // 监控中位数,允许1%的误差
			0.75:  0.01,   // 监控75分位,允许1%的误差
			0.9:   0.01,   // 监控90分位,允许1%的误差
			0.99:  0.001,  // 监控99分位,允许0.1%的误差
			0.999: 0.0001, // 监控99.9分位,允许0.01%的误差
		},
	}, labels)

	// 注册指标
	prometheus.MustRegister(vector)

	// 返回Gin中间件函数
	return func(ctx *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		// 使用defer确保在请求结束时记录耗时
		defer func() {
			// 计算请求耗时(毫秒)
			duration := time.Since(start).Milliseconds()
			// 获取请求方法
			method := ctx.Request.Method
			// 获取请求路径模式
			pattern := ctx.FullPath()
			// 获取响应状态码
			status := ctx.Writer.Status()
			// 记录观测值
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration))
		}()
		// 继续处理请求
		ctx.Next()
	}
}

// BuildActiveRequests 构建监控活跃请求数的中间件
// 使用Gauge类型的指标,用于记录当前正在处理的HTTP请求数量
// 每个新请求进来时计数器+1,请求结束时计数器-1
// 可以反映系统的实时负载情况
func (b *Builder) BuildActiveRequests() gin.HandlerFunc {
	// 创建一个Gauge类型的指标
	// Gauge适合用于记录可以增加和减少的指标值
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		// Namespace/Subsystem/Name 组成指标的完整名称
		// 例如: geektime_daming_webook_gin_http_active_req
		// 注意:指标名称中只能包含字母、数字和下划线
		Namespace: b.NameSpace,            // 指标的命名空间,例如:geektime_daming
		Subsystem: b.Subsystem,            // 指标的子系统,例如:webook
		Name:      b.Name + "_active_req", // 指标名称,例如:gin_http_active_req

		// ConstLabels用于添加静态标签
		// 这里添加instance_id用于区分不同的服务实例
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})

	// 注册指标到Prometheus的默认注册表
	// 如果注册失败会panic
	prometheus.MustRegister(gauge)

	// 返回Gin中间件函数
	return func(ctx *gin.Context) {
		// 新请求到来时将gauge值+1
		gauge.Inc()

		// 使用defer确保在请求处理完成后执行gauge值-1
		// 无论请求是否成功都会执行
		defer gauge.Dec()

		// 调用Next()继续处理请求
		// 等待其他中间件和最终的请求处理程序执行完成
		ctx.Next()
	}
}
