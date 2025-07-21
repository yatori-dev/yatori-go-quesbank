package dbutils

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
	"yatori-go-quesbank/ques-core/entity"
)

func DBClient(dbPath string) (*gorm.DB, error) {
	// 判断文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, errors.New("数据库文件不存在，不打开数据库")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
	return db, nil
}
