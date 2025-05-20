package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Question struct {
	Hash    string   //
	Type    string   //题目类型
	Content string   // 题目内容
	Options []string //如果是选择题那么会返回这个选项
	Answer  []string
	Json    string
}

func main() {

	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		var question Question
		if err := c.ShouldBindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{
			"type": question.Type, //一般原模原样返回
			"answers": []string{
				question.Options[0],
			},
		})
		return
	})
	router.Run(":8083")
}
