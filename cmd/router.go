package cmd

import api_server "yatori-go-quesbank/cmd/api-server"

func (router Group) QuestionRouter() {
	router.POST("/", api_server.QuestionApi)
}
