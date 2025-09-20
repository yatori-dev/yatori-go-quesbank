package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
	"yatori-go-quesbank/global"
	"yatori-go-quesbank/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"yatori-go-quesbank/config"
	"yatori-go-quesbank/ques-core/entity"
)

// 初始化
func BankInit() {
	LogInit()
	fmt.Println(config.YatoriLogo())                                   //打印LOGO
	ConfigInit()                                                       //初始化配置文件读取
	DBInit(global.GlobalConfig.Setting.BasicSetting.DefaultSqlitePath) //初始Sqlite本地数据库
	init := ServerInit()                                               //初始化服务器
	init.Run(":8083")
}

// 日志系统初始化
func LogInit() {
	// 设置日志输出到终端
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 设置最低日志级别为 DEBUG
	})

	logger := slog.New(handler)
	slog.SetDefault(logger) // 设置为默认日志器（可选）
}

// 配置文件系统初始化
func ConfigInit() {
	global.GlobalConfig = config.ReadConfig("./config.yaml")
}

// 数据库初始化
func DBInit(path string) {
	utils.PathExistForCreate(filepath.Dir(path)) //检测目录是否存在，不存在则创建
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
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
	global.GlobalSqliteMap[path] = db
}

type Group struct {
	*gin.RouterGroup
}

// 初始化gin
func ServerInit() *gin.Engine {
	router := gin.Default()
	apiGroup := router.Group("")
	routerGroup := Group{apiGroup}
	routerGroup.QuestionRouter()
	return router
}
