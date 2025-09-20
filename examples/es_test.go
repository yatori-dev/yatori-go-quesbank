package examples

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"yatori-go-quesbank/ques-core/entity"
	questionbank "yatori-go-quesbank/ques-core/localq"

	es9 "github.com/elastic/go-elasticsearch/v9"
)

func TestEsSearch(t *testing.T) {
	cfg := es9.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "ymtberZQNjsomxwduLmv", // 启动日志里生成的密码
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client, err := es9.NewTypedClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	//创建索引
	//createIndex(client, "questions1")
	//将题目插入对应索引
	//indexDocument(client, "questions1", entity.DataQuestion{
	//	Question: entity.Question{
	//		Content: "测试问题1",
	//	},
	//	RightStatus: 0,
	//})
	//查询问题
	questionbank.EsQuestQuestionForContentMachOne(client, "quesbank1", entity.EsQuestion{
		Question: entity.Question{
			Content: "这是问题",
		},
	})
}

// 创建索引
func createIndex(client *es9.TypedClient, indexName string) {
	resp, err := client.Indices.Create(indexName).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

// 创建一条document并添加到对应索引中
func indexDocument(client *es9.TypedClient, indexName string, dataQuestion entity.DataQuestion) {
	do, err := client.Index(indexName).
		Document(dataQuestion).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(do)
}
