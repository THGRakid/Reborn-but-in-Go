package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// DB 数据库链接单例
// 创建了一个名为 DB 的变量，用于在代码中保存一个指向 GORM v2 版本的数据库连接对象的指针
var (
	DB *gorm.DB
)

// 数据库连接的初始化。初始化函数会在程序启动时自动执行。
func init() {
	// 连接数据库
	dsn := "Reborn_root:Go123456@tcp(sh-cynosdbmysql-grp-ni8vyxi2.sql.tencentcdb.com:24351)/reborn_but_in_go"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	})
	if err != nil {
		fmt.Println("数据库连接错误:", err)
		return
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("连接池初始化失败:", err)
		return
	}
	//默认不加复数
	sqlDB.SingularTable(true)
	//设置连接池
	//空闲
	sqlDB.DB().SetMaxIdleConns(20)
	//打开
	sqlDB.DB().SetMaxOpenConns(100)
	//超时
	sqlDB.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db // 将数据库连接保存到全局变量
}

/*
// Database 在中间件中初始化mysql链接

func Database(connString string) {
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
*/
