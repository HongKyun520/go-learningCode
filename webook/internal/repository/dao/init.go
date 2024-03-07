package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {

	// 自动建表,go不存在根据注解扫描的涉及
	return db.AutoMigrate(&User{})

}
