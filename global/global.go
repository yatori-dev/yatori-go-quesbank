package global

import (
	"gorm.io/gorm"
	config2 "yatori-go-quesbank/config"
)

var GlobalDBMap = make(map[string]*gorm.DB) //数据库列表

var GlobalConfig config2.JSONDataForConfig //配置文件
