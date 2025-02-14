package ioc

import (
	"GoInAction/webook/internal/repository/dao"
	"GoInAction/webook/pkg/gormx"

	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func InitDB() *gorm.DB {
	// 设置dsn默认值
	viper.Set("db.dsn", "root:Kyun1024!@tcp(localhost:3306)/webook?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")

	// 获取dsn
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	var c Config
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(err)
	}
	// mysql连接地址 这里只是本地环境的连法
	db, err := gorm.Open(mysql.Open(c.DSN))

	// mysql连接地址，这里是k8s内部的连接法，pod与pod之间通过service的name和port进行连接
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-live-mysql:11309)/webook"))
	if err != nil {
		// panic表示该goroutine直接结束。只在main函数的初始化过程中，使用panic
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	// 接入prometheus统计执行线程数
	// 使用 Prometheus 插件来监控数据库性能
	// DBName: 设置数据库名称为 "webook"
	// RefreshInterval: 设置刷新间隔为 15 秒
	// MetricsCollector: 配置要收集的指标
	//   - MySQL.VariableNames: 收集 thread_running 指标,用于监控当前运行的线程数
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))

	if err != nil {
		panic(err)
	}

	// 接入prometheus统计执行耗时
	err = initPrometheusMetrics(db)
	if err != nil {
		panic(err)
	}

	// 初始化表结构
	err = dao.InitTables(db)

	if err != nil {
		panic(err)
	}

	return db
}

// initPrometheusMetrics 初始化 Prometheus 指标监控
func initPrometheusMetrics(db *gorm.DB) error {
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "geektime_daming",
		Subsystem: "webook",
		Name:      "gorm_db",
		Help:      "统计GORM执行耗时",
		ConstLabels: map[string]string{
			"instance_id": "my_instance_id",
		},
		Objectives: map[float64]float64{0.5: 0.01, 0.75: 0.01, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	})

	err := db.Callback().Create().Before("*").Register("prometheus_gorm_before_create", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Create().After("*").Register("prometheus_gorm_after_create", cb.After("create"))
	if err != nil {
		return err
	}

	err = db.Callback().Query().Before("*").Register("prometheus_gorm_before_query", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Query().After("*").Register("prometheus_gorm_after_query", cb.After("query"))
	if err != nil {
		return err
	}

	err = db.Callback().Raw().Before("*").Register("prometheus_gorm_before_raw", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Raw().After("*").Register("prometheus_gorm_after_raw", cb.After("raw"))
	if err != nil {
		return err
	}

	err = db.Callback().Update().Before("*").Register("prometheus_gorm_before_update", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Update().After("*").Register("prometheus_gorm_after_update", cb.After("update"))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().Before("*").Register("prometheus_gorm_before_delete", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Delete().After("*").Register("prometheus_gorm_after_delete", cb.After("delete"))
	if err != nil {
		return err
	}

	err = db.Callback().Row().Before("*").Register("prometheus_gorm_before_row", cb.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Row().After("*").Register("prometheus_gorm_after_row", cb.After("row"))
	if err != nil {
		return err
	}
	return nil
}
