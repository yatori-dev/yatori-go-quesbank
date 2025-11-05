package yanxi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"
	"yatori-go-quesbank/utils/qutils"

	"github.com/thedevsaddam/gojsonq"
)

// [{"name":"言溪题库","homepage":"https://tk.enncy.cn/","url":"https://tk.enncy.cn/query","method":"get","type":"GM_xmlhttpRequest","contentType":"json","data":{"token":"9e20541d49204bf0813a76e6f3bfdc7e","title":"${title}","options":"${options}","type":"${type}"},"handler":"return (res)=>res.code === 0 ? [res.data.answer, undefined] : [res.data.question,res.data.answer]"},{"name":"网课小工具题库（GO题）","homepage":"https://cx.icodef.com/","url":"https://cx.icodef.com/wyn-nb?v=4","method":"post","type":"GM_xmlhttpRequest","data":{"question":"${title}"},"headers":{"Content-Type":"application/x-www-form-urlencoded","Authorization":""},"handler":"return  (res)=> res.code === 1 ? [undefined,res.data] : [res.msg,undefined]"}]
func questionRequest(token string, title string, qType qtype.QType) string {
	urlStr := "https://tk.enncy.cn/query"
	method := "GET"
	params := url.Values{}
	params.Add("token", token)
	params.Add("title", title)
	params.Add("type", qType.YanXiString())
	client := &http.Client{}
	//slog.Debug("言溪题库请求参数编码--> ", params.Encode())
	req, err := http.NewRequest(method, urlStr+"?"+params.Encode(), nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return string(body)
}

func maxArray(arrs ...[]string) []string {
	if len(arrs) == 0 {
		return nil
	}
	maxArr := arrs[0]
	for _, arr := range arrs[1:] {
		if len(arr) > len(maxArr) {
			maxArr = arr
		}
	}
	return maxArr
}

// {"code":0,"data":{"question":"无","answer":"题目不能为空"},"message":"请求失败"}
// {"code":1,"data":{"question":"测试题目","answer":"用于检验学员受训后知识、技能以及绩效状况的一系列问题或评价方法。","times":98},"message":"请求成功"}
func Request(token string, question entity.Question) *entity.Question {
	resContent := qutils.RemoveLeadingLabel(question.Content)
	jsonStr := questionRequest(token, resContent, qtype.Index(question.Type))
	//jsonStr := QuestionRequest(token, question.Content, question.Type)
	json := gojsonq.New().JSONString(jsonStr)
	if int(json.Find("code").(float64)) != 1 {
		return nil
	}
	//fmt.Println(json)
	//fmt.Println(gojsonq.New().JSONString(jsonStr).Find("data.answer"))
	if gojsonq.New().JSONString(jsonStr).Find("data.answer") == nil {
		return nil
	}
	//第一种分割类型
	answer1 := strings.Split(gojsonq.New().JSONString(jsonStr).Find("data.answer").(string), "#")
	//去空
	answer1 = func(v []string) []string {
		res := []string{}
		for _, answer := range v {
			if answer != "" {
				res = append(res, answer)
			}

		}
		return res
	}(answer1)

	//第二种分割类型
	answer2 := strings.Split(gojsonq.New().JSONString(jsonStr).Find("data.answer").(string), "\n")
	//去空
	answer2 = func(v []string) []string {
		res := []string{}
		for _, answer := range v {
			if answer != "" {
				res = append(res, answer)
			}

		}
		return res
	}(answer2)

	//第三种分割类型
	answer3 := strings.Split(gojsonq.New().JSONString(jsonStr).Find("data.answer").(string), "---")
	//去空
	answer3 = func(v []string) []string {
		res := []string{}
		for _, answer := range v {
			if answer != "" {
				res = append(res, answer)
			}

		}
		return res
	}(answer3)

	// 第二种回复类型
	answer4 := strings.Split(strings.ReplaceAll(strings.ReplaceAll(gojsonq.New().JSONString(jsonStr).Find("data.answer").(string), "[", ""), "]", ""), ",")
	//去空
	answer4 = func(v []string) []string {
		res := []string{}
		for _, answer := range v {
			if answer != "" {
				res = append(res, answer)
			}

		}
		return res
	}(answer4)
	//赋值可以分组最大的那个作为答案
	question.Answers = maxArray(answer1, answer2, answer3, answer4)

	//检测是否为选项字母答案，如果是，则转换
	for i, option := range question.Answers {
		if option == "A" && len(question.Options) >= 1 {
			question.Answers[i] = question.Options[0]
		} else if option == "B" && len(question.Options) >= 2 {
			question.Answers[i] = question.Options[1]
		} else if option == "C" && len(question.Options) >= 3 {
			question.Answers[i] = question.Options[2]
		} else if option == "D" && len(question.Options) >= 4 {
			question.Answers[i] = question.Options[3]
		} else if option == "E" && len(question.Options) >= 5 {
			question.Answers[i] = question.Options[4]
		} else if option == "F" && len(question.Options) >= 6 {
			question.Answers[i] = question.Options[5]
		}
	}

	//题目也赋值一样
	queryContent := gojsonq.New().JSONString(jsonStr).Find("data.question")
	if queryContent != nil && queryContent.(string) != "" {
		question.Content = queryContent.(string)
	}

	return &question
}
