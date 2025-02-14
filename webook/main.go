package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

// main方法需要自己手写很多东西

// go项目需要在main方法中手动做很多事情
// 初始化数据源、甚至初始化普通类、手动注册接口路由、手动注册middleware(过滤器、拦截器)

// 应用启动入口
func main() {

	// 初始化配置文件
	initViperWatch()

	// 初始化日志
	initLogger()

	// 初始化观测
	initPrometheus()

	//初始化服务器（使用wire）
	server := InitWebServer()

	server.Run(":8080")

}

// initLogger 初始化日志组件
// 使用zap作为日志库
// 开发环境使用NewDevelopment()创建logger
// 如果创建失败则panic
// 最后将创建的logger替换到全局logger
func initLogger() {
	// 创建开发环境的logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// 替换全局的logger
	zap.ReplaceGlobals(logger)
}

func initViperWatch() {

	// 可选参数，默认值为config/dev.yaml。实现在不同环境下读取不同的配置文件
	cflag := pflag.String("config", "config/dev.yaml", "配置文件路径")

	// 解析
	pflag.Parse()

	log.Println("启动参数", *cflag)

	// 设置配置文件名
	viper.SetConfigFile(*cflag)

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件变化", viper.GetString("test.name"))
	})

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	log.Println(viper.GetString("test.name"))
}

func initViper() {

	// 设置配置文件名
	viper.SetConfigName("dev")

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 当前目录下的config目录
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	log.Println(viper.GetString("test.name"))
}

func initViperV1() {

	// 可选参数，默认值为config/dev.yaml。实现在不同环境下读取不同的配置文件
	cflag := pflag.String("config", "config/dev.yaml", "配置文件路径")

	// 解析
	pflag.Parse()

	log.Println("启动参数", *cflag)

	// 设置配置文件名
	viper.SetConfigFile(*cflag)

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	log.Println(viper.GetString("test.name"))

}

func initViperRemote() {

	// 设置etcd地址
	// 第一个参数是远程配置类型,这里使用etcd3
	// 第二个参数是etcd的访问地址,这里使用本地地址和端口12379
	// 第三个参数是配置在etcd中的key路径,这里使用/webook
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 读取配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件变化", viper.GetString("test.name"))
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	log.Println(viper.GetString("test.name"))

	// 监听配置文件变化
	go func() {
		for {
			err := viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("配置文件变化", viper.GetString("test.name"))
			time.Sleep(time.Second * 10)
		}
	}()

}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()
}
