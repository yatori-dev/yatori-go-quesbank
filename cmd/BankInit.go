package cmd

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"time"
	api_server "yatori-go-quesbank/cmd/api-server"
	"yatori-go-quesbank/global"

	"yatori-go-quesbank/config"
	"yatori-go-quesbank/ques-core/entity"
)

// 初始化
func BankInit() {
	LogInit()
	fmt.Println(config.YatoriLogo()) //打印LOGO
	DBInit()                         //初始化数据库
	ServerInit()                     //初始化服务器
}
func LogInit() {
	// 设置日志输出到终端
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 设置最低日志级别为 DEBUG
	})

	logger := slog.New(handler)
	slog.SetDefault(logger) // 设置为默认日志器（可选）
}

// 数据库初始化
func DBInit() {
	db, err := gorm.Open(sqlite.Open("localq.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&entity.DataQuestion{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB() //数据库连接池
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.GlobalDB = db
}

// 初始化gin
func ServerInit() {
	api_server.QuestionApi()
}
