package ques_core

import (
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"yatori-go-quesbank/config"
	"yatori-go-quesbank/global"
	"yatori-go-quesbank/ques-core/aiq"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"
	"yatori-go-quesbank/ques-core/externalq/yanxi"
	questionbank "yatori-go-quesbank/ques-core/localq"
	"yatori-go-quesbank/utils/dbutils"

	es9 "github.com/elastic/go-elasticsearch/v9"
)

// 执行逻辑
func Research(anSet []config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	//循环获取配置文件中的题库设置
	for _, v := range anSet {
		var que *entity.DTOQuestion
		switch v.AnswerType {
		case "SQLITE": //sqlite
			que = sqliteResearch(v, question)
		case "ES":
			que = esResearch(v, question)
		case "AI": //AI
			que = aiResearch(v, question)
		case "EXTERNAL": //外部三方题库
			que = externResearch(v, question)
		}
		if que != nil {
			//自动缓存逻辑
			AutoCaches(v.CacheTargetList, anSet, que.Question)
			return que
		}
	}
	return nil
}

// 自动查询
func AutoResearch(question entity.Question) *entity.DTOQuestion {
	return Research(global.GlobalConfig.Setting.AnswerSetting, question)
}

// 自动缓存
func AutoCaches(cacheTargetList []string, anSet []config.AnswerSetting, question entity.Question) {
	for _, v := range cacheTargetList {
		configData := config.GetAnswerConfigForLabel(anSet, v)
		if configData == nil {
			continue
		}
		//根据类型缓存
		switch configData.AnswerType {
		case "SQLITE":
			//懒加载Sqlite
			LoadSqlite(*configData)
			//缓存sqlite
			err := questionbank.InsertIfNot(global.GlobalSqliteMap[global.GlobalConfig.Setting.BasicSetting.DefaultSqlitePath], &entity.DataQuestion{Question: question})
			if err != nil {
				log.Println(err)
			}
		case "ES":
			//懒加载ES
			LoadEs(*configData)
			//缓存es
			err := questionbank.EsInsertIfNot(global.GlobalESMap[configData.EsUrl], configData.EsIndex, entity.EsQuestion{Question: question})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// 本地全部查询
func LocalAllResearch() (result []entity.DTOQuestion) {
	db := global.GlobalSqliteMap[global.GlobalConfig.Setting.BasicSetting.DefaultSqlitePath]
	allQuestion := questionbank.SelectsAllQuestion(db)
	for _, que := range allQuestion {
		result = append(result, entity.DTOQuestion{Question: que.Question, Replier: "Sqlite", ReplyType: "LOCAL"})
	}
	return result
}

// 本地类型查询
func LocalTypeResearch(qtype qtype.QType) (result []entity.DTOQuestion) {
	db := global.GlobalSqliteMap[global.GlobalConfig.Setting.BasicSetting.DefaultSqlitePath]
	forTypeQues := questionbank.SelectsForType(db, qtype)
	for _, que := range forTypeQues {
		result = append(result, entity.DTOQuestion{Question: que.Question, Replier: "Sqlite", ReplyType: "LOCAL"})
	}
	return
}

// Sqlite搜索
func sqliteResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	//探测是否本地缓存库有
	//resultData := questionbank.SelectForTypeAndContent(global.GlobalDB, &entity.DataQuestion{Question: question})
	var resultData *entity.DataQuestion
	if anSet.SqlitePath == "" {
		//默认的本地数据库查询
		resultData = questionbank.SelectForTypeAndLikeContent1_4(global.GlobalSqliteMap[global.GlobalConfig.Setting.BasicSetting.DefaultSqlitePath], &entity.DataQuestion{Question: question})
	} else {
		//懒加载Sqlite
		LoadSqlite(anSet)
		resultData = questionbank.SelectForTypeAndLikeContent1_4(global.GlobalSqliteMap[anSet.SqlitePath], &entity.DataQuestion{Question: question})
	}

	if resultData == nil {
		return nil
	}
	return &entity.DTOQuestion{Question: resultData.Question, Replier: anSet.AnswerLabel, ReplyType: "LOCAL"}
}

// 加载Sqlite
func LoadSqlite(anSet config.AnswerSetting) {
	//懒加载数据库
	if global.GlobalSqliteMap[anSet.SqlitePath] == nil {
		db, err := dbutils.DBClient(anSet.SqlitePath)
		if err != nil {
			log.Println(err.Error())
		}
		global.GlobalSqliteMap[anSet.SqlitePath] = db
	}
}

// ElasticSearch搜索
func esResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	//如果该ES没有填写url，那么直接采用默认的
	if anSet.EsUrl == "" {
		anSet.EsUrl = global.GlobalConfig.Setting.BasicSetting.DefaultEsUrl
	}
	//懒加载ES
	LoadEs(anSet)
	//获取es对象
	esClient := global.GlobalESMap[global.GlobalConfig.Setting.BasicSetting.DefaultEsUrl]
	questionData := questionbank.EsQuestQuestionForContentMachOne(esClient, anSet.EsIndex, entity.EsQuestion{Question: question})
	if questionData == nil {
		return nil
	}
	return &entity.DTOQuestion{Question: *questionData, Replier: anSet.AnswerLabel, ReplyType: anSet.AnswerType}
}

// 加载ES
func LoadEs(anSet config.AnswerSetting) {
	//懒加载ES
	if global.GlobalESMap[global.GlobalConfig.Setting.BasicSetting.DefaultEsUrl] == nil {
		cfg := es9.Config{
			Addresses: []string{
				global.GlobalConfig.Setting.BasicSetting.DefaultEsUrl,
			},
			Username: anSet.EsUsername,
			Password: anSet.EsPassword, // 启动日志里生成的密码
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: anSet.EsSkipVerify},
			},
		}
		client, err := es9.NewTypedClient(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
		global.GlobalESMap[global.GlobalConfig.Setting.BasicSetting.DefaultEsUrl] = client
		//查询是否有对应的索引，如果没有则创建
		questionbank.EsQuestIndexOrNotForCreate(client, anSet.EsIndex)
	}
}

// AI搜索
func aiResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	aiAnswer, err := aiq.AggregationAIApi(anSet.AiUrl, anSet.AiModel, anSet.AiType, aiq.AIChatMessages{}, anSet.APIKEY)
	slog.Debug(aiAnswer, err)
	return nil
}

// 外部题库搜索
func externResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	switch anSet.ExType {
	case "YANXI":
		//使用言溪题库
		result := yanxi.Request(anSet.ExToken, question)
		if result != nil {
			return &entity.DTOQuestion{Question: *result, Replier: anSet.AnswerLabel, ReplyType: "EXTERNAL"}
		}
	}

	return nil
}
