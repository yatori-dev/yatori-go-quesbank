package yanxi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"

	"github.com/thedevsaddam/gojsonq"
)

// [{"name":"言溪题库","homepage":"https://tk.enncy.cn/","url":"https://tk.enncy.cn/query","method":"get","type":"GM_xmlhttpRequest","contentType":"json","data":{"token":"9e20541d49204bf0813a76e6f3bfdc7e","title":"${title}","options":"${options}","type":"${type}"},"handler":"return (res)=>res.code === 0 ? [res.data.answer, undefined] : [res.data.question,res.data.answer]"},{"name":"网课小工具题库（GO题）","homepage":"https://cx.icodef.com/","url":"https://cx.icodef.com/wyn-nb?v=4","method":"post","type":"GM_xmlhttpRequest","data":{"question":"${title}"},"headers":{"Content-Type":"application/x-www-form-urlencoded","Authorization":""},"handler":"return  (res)=> res.code === 1 ? [undefined,res.data] : [res.msg,undefined]"}]
func questionRequest(token string, title string, qType qtype.QType) string {

	urlStr := "https://max.tlicf.com/Interface/xxt/?key=" + token + "&question=" + url.QueryEscape(title)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, urlStr, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "maxq.tlicf.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return string(body)
}

// {"code":0,"data":{"question":"无","answer":"题目不能为空"},"message":"请求失败"}
// {
// "code": 1,
// "data": "活期储蓄###整存整取###定活两便###通知存款"
// }
func Request(token string, question entity.Question) *entity.Question {
	jsonStr := questionRequest(token, question.Content, qtype.Index(question.Type))
	//jsonStr := QuestionRequest(token, question.Content, question.Type)
	json := gojsonq.New().JSONString(jsonStr)
	if int(json.Find("code").(float64)) != 1 {
		return nil
	}
	fmt.Println(json)
	fmt.Println(gojsonq.New().JSONString(jsonStr).Find("data"))
	if gojsonq.New().JSONString(jsonStr).Find("data") == nil {
		return nil
	}
	//第一种回复类型
	answer1 := strings.Split(gojsonq.New().JSONString(jsonStr).Find("data").(string), "###")
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
	question.Answers = answer1

	return &question
}
