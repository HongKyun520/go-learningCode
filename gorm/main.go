package main

import (
	"fmt"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type Product struct {
	ID    string `gorm:"primaryKey"`
	Code  string `gorm:"column:code"`
	Price uint   `gorm:"column:price"`
}

// 使用gorm建议都是使用指针
func main() {
	// 根据配置创建数据库实例
	// 使用sqlLite
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	//	// 只输出sql语句 但不执行
	//	DryRun: true,
	//})

	// 使用MySQL
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/test"))
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/your_db"))

	db = db.Debug()

	if err != nil {
		panic("failed to connect database")
	}

	db = db.Debug()

	// 建表
	// 表存在则改为一样
	db.AutoMigrate(&Product{})

	// Create  insert
	db.Create(&Product{ID: "1", Code: "D42", Price: 100})
	db.Create(&Product{ID: "2", Code: "D43", Price: 101})

	// Read   select
	var product Product
	//db.First(&product, 1) // 根据整型主键查找

	db.First(&product, "code = ?", "D43") // 查找 code 字段值为 D42 的记录
	fmt.Printf("Product: %+v\n", product)

	// Update - 将 product 的 price 更新为 200
	//db.Model(&product).Update("Price", 200)

	// 将 product 的 code 更新为 "D43"  orm框架不做value类型的校验
	//db.Model(&product).Update("Code", "D43")

	// Update - 更新多个字段  传入结构体
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段

	// Update - 更新多个字段  传入Map
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	//db.Delete(&product, 1)

}
