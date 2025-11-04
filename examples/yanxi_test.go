package examples

import (
	"fmt"
	"testing"
	yanxi "yatori-go-quesbank/ques-core/externalq/zaodianshui"
)

// 测试言溪题库请求
func TestRequestApi(t *testing.T) {
	fmt.Println(yanxi.RemoveLeadingLabel("下列有关国债逆回购的说法正确的是（）"))
	//yanxi. Request("", "撒打算大萨达撒大声地", qtype.SingleChoice)
}
