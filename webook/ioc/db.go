package ioc

import (
	"GoInAction/webook/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// mysql连接地址 这里只是本地环境的连法
	db, err := gorm.Open(mysql.Open("root@tcp(localhost:3306)/webook"))

	// mysql连接地址，这里是k8s内部的连接法，pod与pod之间通过service的name和port进行连接
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-live-mysql:11309)/webook"))
	if err != nil {
		// panic表示该goroutine直接结束。只在main函数的初始化过程中，使用panic
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	// 初始化表结构
	err = dao.InitTable(db)

	if err != nil {
		panic(err)
	}

	return db
}
