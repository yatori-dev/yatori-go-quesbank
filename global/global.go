package global

import (
	config2 "yatori-go-quesbank/config"

	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

var GlobalSqliteMap = make(map[string]*gorm.DB)               //Sqlite数据库列表
var GlobalESMap = make(map[string]*elasticsearch.TypedClient) //ElasticSearch列表

var GlobalConfig config2.JSONDataForConfig //配置文件
