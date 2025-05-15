package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Question struct {
	Hash    string
	Type    string
	Content string
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
			"type": "单选",
			"answers": []string{
				"测试",
			},
		})
		return
	})
	router.Run(":8083")
}
