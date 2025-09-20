package questionbank

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"yatori-go-quesbank/ques-core/entity"

	es9 "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// 根据题目查询
func EsQuestQuestionForContentMachOne(client *es9.TypedClient, indexName string, esQuestion entity.EsQuestion) *entity.Question {
	resp, err := client.Search().
		Index(indexName).
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"content": {Query: esQuestion.Content, Fuzziness: "AUTO"},
			},
		}).Size(1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		return nil
	}
	fmt.Printf("total: %d\n", resp.Hits.Total.Value)
	if resp.Hits.Total.Value == 0 {
		return nil
	}

	question := &entity.Question{}
	if err1 := json.Unmarshal(resp.Hits.Hits[0].Source_, question); err1 != nil {
		fmt.Printf("unmarshal failed: %v\n", err)
		return nil
	}
	return question
}

// 查询对应ES是否有索引，没有则直接创建
func EsQuestIndexOrNotForCreate(client *es9.TypedClient, indexName string) {
	exisit, err := client.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		fmt.Printf("check index exists failed, err:%v\n", err)
	}
	if !exisit {
		_, err1 := client.Indices.Create(indexName).Do(context.Background())
		if err1 != nil {
			panic(err1)
		}
	}
}

// 创建一条document并添加到对应索引中
func EsInsert(client *es9.TypedClient, indexName string, esQuestion entity.EsQuestion) error {
	do, err := client.Index(indexName).
		Document(esQuestion).
		Do(context.Background())
	if err != nil {
		panic(err)
		return err
	}
	if do.Result.String() != "created" {
		return errors.New("插入ES失败，返回信息：" + do.Result.String())
	}

	return nil
}

// 如果没有对应题目，则直接插入
func EsInsertIfNot(client *es9.TypedClient, indexName string, esQuestion entity.EsQuestion) error {
	resp, err := client.Search().
		Index(indexName).
		Query(&types.Query{
			Term: map[string]types.TermQuery{
				"content.keyword": { // 注意用 .keyword
					Value: esQuestion.Content,
				},
			},
		}).
		Size(1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
	}
	//如果没有一模一样的，则插入
	if resp.Hits.Total.Value == 0 {
		err1 := EsInsert(client, indexName, esQuestion)
		if err1 != nil {
			return err1
		}
	}
	return nil
}
