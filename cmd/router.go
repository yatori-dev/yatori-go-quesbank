package cmd

import api_server "yatori-go-quesbank/cmd/api-server"

func (router Group) QuestionRouter() {
	var localApi api_server.LocalQuestionApi
	router.POST("/", api_server.QuestionApi)
	router.GET("/all_question", localApi.AllQuestionApi)
	router.GET("/all_question/:type", localApi.SelectsTypeApi)
}
