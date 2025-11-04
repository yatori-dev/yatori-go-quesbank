package examples

import (
	"fmt"
	"testing"
	"yatori-go-quesbank/ques-core/entity"
	zds "yatori-go-quesbank/ques-core/externalq/zaodianshui"
)

// 测试言溪题库请求
func TestZDSRequestApi(t *testing.T) {
	question := entity.Question{
		Content: "药用植物是指具有____,____功能的植物。",
	}
	request := zds.Request("", question)
	fmt.Println(request)
	//fmt.Println(zds.RemoveLeadingLabel("下列有关国债逆回购的说法正确的是（）"))
	//yanxi. Request("", "撒打算大萨达撒大声地", qtype.SingleChoice)
}
