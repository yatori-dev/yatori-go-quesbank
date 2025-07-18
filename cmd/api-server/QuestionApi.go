package api_server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yatori-go-quesbank/global"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/replier"
	"yatori-go-quesbank/ques-core/externalq/yanxi"
	"yatori-go-quesbank/ques-core/localq"
)

// 题目请求Api
func QuestionApi() {
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		var question entity.Question
		if err := c.ShouldBindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		//探测是否本地缓存库用
		resultData := questionbank.SelectForTypeAndContent(global.GlobalDB, &entity.DataQuestion{Question: question})
		if resultData != nil {
			c.JSON(http.StatusOK, entity.ResultQuestion{
				Question: resultData.Question,
				Replier:  replier.LOCAL.String(),
				Msg:      "查询成功",
				Code:     200,
			})
			return
		}

		//使用言溪题库
		result := yanxi.Request("", question)
		if result != nil {
			c.JSON(http.StatusOK, entity.ResultQuestion{
				Question: *result,
				Replier:  replier.YANXI.String(),
				Msg:      "查询成功",
				Code:     200,
			})
			err := questionbank.InsertIfNot(global.GlobalDB, &entity.DataQuestion{Question: *result})
			if err != nil {
				log.Println(err)
			}
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
