package dao

import "gorm.io/gorm"

// 初始化简表过程
func InitTable(db *gorm.DB) error {

	// 根据model自动建表,go不存在根据注解扫描的涉及
	return db.AutoMigrate(&User{})

}
