package api_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ques_core "yatori-go-quesbank/ques-core"
	"yatori-go-quesbank/ques-core/entity"
)

// 题目请求Api
func QuestionApi() {
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		var question entity.Question
		if err := c.ShouldBindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		//查询题目
		result := ques_core.AutoResearch(question)
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
	})
	router.Run(":8083")
}
