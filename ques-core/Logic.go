package ques_core

import (
	"log"
	"log/slog"
	"yatori-go-quesbank/config"
	"yatori-go-quesbank/global"
	"yatori-go-quesbank/ques-core/aiq"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"
	"yatori-go-quesbank/ques-core/externalq/yanxi"
	questionbank "yatori-go-quesbank/ques-core/localq"
	"yatori-go-quesbank/utils/dbutils"
)

// 执行逻辑
func Research(anSet []config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	for _, v := range anSet {
		switch v.AnswerType {
		case "LOCAL":
			que := localResearch(v, question)
			if que != nil {
				return que
			}
		case "AI":
			que := aiResearch(v, question)
			if que != nil {
				return que
			}
		case "EXTERNAL":
			que := externResearch(v, question)
			if que != nil {
				return que
			}
		}
	}
	return nil
}

// 自动查询
func AutoResearch(question entity.Question) *entity.DTOQuestion {
	return Research(global.GlobalConfig.Setting.AnswerSetting, question)
}

func LocalAllResearch() (result []entity.DTOQuestion) {
	db := global.GlobalDBMap[global.GlobalConfig.Setting.BasicSetting.DefaultDBPath]
	allQuestion := questionbank.SelectsAllQuestion(db)
	for _, que := range allQuestion {
		result = append(result, entity.DTOQuestion{Question: que.Question, Replier: "LOCAL", ReplyType: "LOCAL"})
	}
	return result
}

func LocalTypeResearch(qtype qtype.QType) (result []entity.DTOQuestion) {
	db := global.GlobalDBMap[global.GlobalConfig.Setting.BasicSetting.DefaultDBPath]
	forTypeQues := questionbank.SelectsForType(db, qtype)
	for _, que := range forTypeQues {
		result = append(result, entity.DTOQuestion{Question: que.Question, Replier: "LOCAL", ReplyType: "LOCAL"})
	}
	return
}

// 本地搜索
func localResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	//探测是否本地缓存库用
	//resultData := questionbank.SelectForTypeAndContent(global.GlobalDB, &entity.DataQuestion{Question: question})
	var resultData *entity.DataQuestion
	if anSet.LocalPath == "" {
		resultData = questionbank.SelectForTypeAndLikeContent1_4(global.GlobalDBMap[global.GlobalConfig.Setting.BasicSetting.DefaultDBPath], &entity.DataQuestion{Question: question})
	} else {
		//懒加载数据库
		if global.GlobalDBMap[anSet.LocalPath] == nil {
			db, err := dbutils.DBClient(anSet.LocalPath)
			if err != nil {
				log.Println(err.Error())
			}
			global.GlobalDBMap[anSet.LocalPath] = db
		}
		resultData = questionbank.SelectForTypeAndLikeContent1_4(global.GlobalDBMap[anSet.LocalPath], &entity.DataQuestion{Question: question})
	}

	if resultData == nil {
		return nil
	}
	return &entity.DTOQuestion{Question: resultData.Question, Replier: "LOCAL", ReplyType: "LOCAL"}
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
			//缓存本地
			if anSet.AutoCache == 1 {
				err := questionbank.InsertIfNot(global.GlobalDBMap[global.GlobalConfig.Setting.BasicSetting.DefaultDBPath], &entity.DataQuestion{Question: *result})
				if err != nil {
					log.Println(err)
				}
			}

			return &entity.DTOQuestion{Question: *result, Replier: "YANXI", ReplyType: "EXTERNAL"}
		}
	}

	return nil
}
