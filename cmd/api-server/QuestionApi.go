package api_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ques_core "yatori-go-quesbank/ques-core"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"
)

type LocalQuestionApi struct{}

func (LocalQuestionApi) AllQuestionApi(c *gin.Context) {
	questions := ques_core.LocalAllResearch()
	if questions != nil {
		c.JSON(http.StatusOK, entity.ListQuestion[[]entity.DTOQuestion]{
			Count:   int64(len(questions)),
			List:    questions,
			Msg:     "查询成功",
			Replier: "LOCAL",
			Code:    200,
		})
	}
	c.JSON(http.StatusOK, entity.ResultQuestion{
		Msg:  "未找到",
		Code: 404,
	})
	return
}

func (LocalQuestionApi) SelectsTypeApi(c *gin.Context) {
	queType := c.Param("type")
	questions := ques_core.LocalTypeResearch(qtype.Index(queType))
	if questions != nil {
		c.JSON(http.StatusOK, entity.ListQuestion[[]entity.DTOQuestion]{
			Count:   int64(len(questions)),
			List:    questions,
			Msg:     "查询成功",
			Replier: "LOCAL",
			Code:    200,
		})
	}
	c.JSON(http.StatusOK, entity.ResultQuestion{
		Msg:  "未找到",
		Code: 404,
	})

	return
}

type QuestionRequest struct {
	Type    qtype.QType        `json:"type"`
	Content string             `json:"content"`
	Options entity.StringArray `json:"options"`
}

// 题目请求Api
func QuestionApi(c *gin.Context) {
	var question QuestionRequest

	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//查询题目
	result := ques_core.AutoResearch(entity.Question{
		Md5:     "",
		Type:    question.Type,
		Content: question.Content,
		Options: question.Options,
		Answers: nil,
	})
	if result != nil {
		c.JSON(http.StatusOK, entity.ResultQuestion{
			Question: result.Question,
			Replier:  result.Replier,
			Msg:      "查询成功",
			Code:     200,
		})
		return
	}
	c.JSON(http.StatusOK, entity.ResultQuestion{
		Msg:  "未找到",
		Code: 404,
	})
	return
}
