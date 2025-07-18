package global

import (
	"gorm.io/gorm"
	config2 "yatori-go-quesbank/config"
)

var GlobalDB *gorm.DB

var GlobalConfig config2.JSONDataForConfig //配置文件
