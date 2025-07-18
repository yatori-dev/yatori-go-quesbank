package ques_core

import (
	"log"
	"log/slog"
	"yatori-go-quesbank/config"
	"yatori-go-quesbank/global"
	"yatori-go-quesbank/ques-core/aiq"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/externalq/yanxi"
	questionbank "yatori-go-quesbank/ques-core/localq"
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

// 本地搜索
func localResearch(anSet config.AnswerSetting, question entity.Question) *entity.DTOQuestion {
	//探测是否本地缓存库用
	//resultData := questionbank.SelectForTypeAndContent(global.GlobalDB, &entity.DataQuestion{Question: question})
	resultData := questionbank.SelectForTypeAndLikeContent1_4(global.GlobalDB, &entity.DataQuestion{Question: question})
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
				err := questionbank.InsertIfNot(global.GlobalDB, &entity.DataQuestion{Question: *result})
				if err != nil {
					log.Println(err)
				}
			}

			return &entity.DTOQuestion{Question: *result, Replier: "YANXI", ReplyType: "EXTERNAL"}
		}
	}

	return nil
}
