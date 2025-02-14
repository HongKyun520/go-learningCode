package gormx

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// Callbacks 结构体用于处理GORM的回调函数,实现Prometheus监控指标的收集
type Callbacks struct {
	// vector 用于记录数据库操作的耗时分布情况
	vector *prometheus.SummaryVec
}

// NewCallbacks 创建一个新的Callbacks实例
// opts: Prometheus指标的配置选项
// 返回: 初始化好的Callbacks实例
func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	// 创建一个SummaryVec,标签包含table(表名)和operation(操作类型)
	vector := prometheus.NewSummaryVec(opts, []string{"table", "operation"})
	// 注册到Prometheus全局注册表
	prometheus.MustRegister(vector)
	return &Callbacks{
		vector: vector,
	}
}

// Before 返回GORM执行SQL前的回调函数
// 返回: 回调函数,用于记录操作开始时间
func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 记录开始时间
		startTime := time.Now()
		// 将开始时间保存到GORM的上下文中
		db.Set("start_time", startTime)
	}
}

// After 返回GORM执行SQL后的回调函数
// typ: 操作类型(如"query"、"create"等)
// 返回: 回调函数,用于计算并记录操作耗时
func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 从上下文中获取开始时间
		startTime, _ := db.Get("start_time")
		start, ok := startTime.(time.Time)
		if !ok {
			return
		}
		// 计算操作耗时(毫秒)
		duration := time.Since(start).Milliseconds()
		// 记录耗时指标,标签包含操作类型和表名
		c.vector.WithLabelValues(typ, db.Statement.Table).Observe(float64(duration))
	}
}
