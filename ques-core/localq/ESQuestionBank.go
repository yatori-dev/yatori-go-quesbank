package questionbank

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"yatori-go-quesbank/ques-core/entity"

	es9 "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// CompareQueryPercent 计算 query 在单条 content 中出现的百分比
func CompareQueryPercent(content string, query string) float64 {
	contentRunes := []rune(content)
	queryRunes := []rune(query)

	count := 0
	for _, qr := range queryRunes {
		for _, cr := range contentRunes {
			if qr == cr {
				count++
				break // 每个 query 字符只算一次
			}
		}
	}
	percent := float64(count) / float64(max(len(queryRunes), len(contentRunes))) * 100
	return percent
}

// 根据题目查询
func EsQuestQuestionForContentMachOne(client *es9.TypedClient, indexName string, esQuestion entity.EsQuestion) *entity.Question {
	queries := []types.Query{
		//匹配内容字段（模糊匹配）
		{
			Match: map[string]types.MatchQuery{
				"content": {Query: esQuestion.Content, Fuzziness: "AUTO"},
			},
		},
		//匹配题目类型
		{
			Term: map[string]types.TermQuery{
				"type.keyword": {Value: esQuestion.Type},
			},
		},
		//判断options中是否包含某个选项值
	}
	//如果是选择题，则再加上选项匹配的可能性
	if (esQuestion.Type == "单选题" || esQuestion.Type == "多选题") && len(esQuestion.Options) > 0 {
		queries = append(queries, types.Query{
			Term: map[string]types.TermQuery{
				"options.keyword": {Value: esQuestion.Options[0]}, // 假设你结构体中有字段 Option
			},
		})
	}
	resp, err := client.Search().
		Index(indexName).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Must: queries,
			},
		}).Size(1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		return nil
	}
	//fmt.Printf("total: %d\n", resp.Hits.Total.Value)
	if resp.Hits.Total.Value == 0 {
		return nil
	}

	question := &entity.Question{}
	if err1 := json.Unmarshal(resp.Hits.Hits[0].Source_, question); err1 != nil {
		fmt.Printf("unmarshal failed: %v\n", err)
		return nil
	}
	//如果查询的题相差太大
	if CompareQueryPercent(question.Content, esQuestion.Content) < 60 {
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

// 分页查询
func EsQueryAll(client *es9.TypedClient, indexName string, pageNum, pageSize int) []entity.Question {
	from := (pageNum - 1) * pageSize
	if from < 0 {
		from = 0
	}

	// 这里只用最“保守”、所有版本几乎都有的字段：
	// Index / From / Size / Query
	req := client.Search().
		Index(indexName).
		From(from).
		Size(pageSize).
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		})

	// 发请求
	res, err := req.Do(context.Background())
	if err != nil {
		log.Printf("es search error: %v", err)
		return nil
	}

	if len(res.Hits.Hits) == 0 {
		return nil
	}

	qs := make([]entity.Question, 0, len(res.Hits.Hits))
	for _, h := range res.Hits.Hits {
		var q entity.Question
		if err := json.Unmarshal(h.Source_, &q); err != nil {
			// 某条解析失败就跳过，不影响其他的
			log.Printf("unmarshal hit error: %v", err)
			continue
		}
		qs = append(qs, q)
	}

	return qs
}
