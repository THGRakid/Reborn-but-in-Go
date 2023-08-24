package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// DB 数据库链接单例
// 创建了一个名为 DB 的变量，用于在代码中保存一个指向 GORM v2 版本的数据库连接对象的指针
var DB *gorm.DB

// 数据库连接的初始化。初始化函数会在程序启动时自动执行。
func init() {

	// MySQL 配置信息
	const userName = "Reborn_root"                                 // 账号
	const password = "Go123456"                                    // 密码
	const host = "sh-cynosdbmysql-grp-ni8vyxi2.sql.tencentcdb.com" // 地址
	const port = 24351                                             // 端口
	const DBName = "reborn_but_in_go"                              // 数据库名称
	const timeOut = "10s"                                          // 连接超时，10秒

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", userName, password, host, port, DBName, timeOut)

	// Open 连接
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

	//空闲
	sqlDB.SetMaxIdleConns(20)
	//打开
	sqlDB.SetMaxOpenConns(100)
	//超时
	sqlDB.SetConnMaxLifetime(time.Second * 30)

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
