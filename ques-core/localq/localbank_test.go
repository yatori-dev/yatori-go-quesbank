package questionbank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"yatori-go-quesbank/ques-core/entity"
)

type CData struct {
	Type    string
	Content string
	Options []string
	Answers []string
}

func TestImportLocalBank(t *testing.T) {
	init, _ := QuestionBankInit()
	// 打开 JSON 文件
	file, err := os.Open("D:\\QQNT\\Downloads\\1.json")
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	var arr = []CData{}
	err = json.Unmarshal(bytes, &arr)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}
	for _, v := range arr {
		fmt.Println(v.Content)
		err := InsertIfNot(init, &entity.DataQuestion{
			Question: entity.Question{
				Type:    v.Type,
				Content: v.Content,
				Options: v.Options,
				Answers: v.Answers,
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
