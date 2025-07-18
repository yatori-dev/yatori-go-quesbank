package examples

import (
	"testing"
	"yatori-go-quesbank/ques-core/entity/qtype"
	"yatori-go-quesbank/ques-core/externalq/yanxi"
)

// 测试言溪题库请求
func TestRequestApi(t *testing.T) {
	yanxi.QuestionRequest("9e20541d49204bf0813a76e6f3bfdc7e", "撒打算大萨达撒大声地", qtype.SingleChoice)
}
