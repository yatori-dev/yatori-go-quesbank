package questionbank

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"yatori-go-quesbank/ques-core/entity"

	"github.com/elastic/go-elasticsearch/v9"
)

type CData struct {
	Type    string
	Content string
	Options []string
	Answers []string
}

func TestImportLocalBank(t *testing.T) {
	init, _ := SqliteQuestionBankInit()
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

type Hit struct {
	Source json.RawMessage `json:"_source"`
}

type SearchResponse struct {
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}

// 导出
func TestExportES(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: "elastic",
		Password: "RnLtD-gnmdOXDgKUS*mM",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("ES init error: %v", err)
	}

	file := "backup.ndjson"
	srcIndex := "questions1"
	//dstIndex := "dst_index"

	// 导出为本地文件
	if err := ExportToFile(client, srcIndex, file); err != nil {
		log.Fatalf("Export error: %v", err)
	}
	fmt.Println("Export to file completed.")

	// 从文件导入到新的 index
	//if err := ImportFromFile(client, dstIndex, file); err != nil {
	//	log.Fatalf("Import error: %v", err)
	//}
	//fmt.Println("Import from file completed.")
}

// 导入
func TestImportES(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: "elastic",
		Password: "RnLtD-gnmdOXDgKUS*mM",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("ES init error: %v", err)
	}

	file := "backup.ndjson"
	dstIndex := "questions1"

	// 从文件导入到新的 index
	if err := ImportFromFile(client, dstIndex, file); err != nil {
		log.Fatalf("Import error: %v", err)
	}
	fmt.Println("Import from file completed.")
}

// 题库去重
func TestReData(t *testing.T) {
	inputFile := "./backup.ndjson"
	outputFile := "./output.ndjson"
	err := DedupLines(inputFile, outputFile)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("去重完成！结果写入:", outputFile)
	}
}

// DedupLines 按行去重并写入新文件（不改变顺序）
func DedupLines(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	seen := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			writer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}

//////////////////////////////////////////////////////
// 导出：使用 scroll → 写入 NDJSON 文件
//////////////////////////////////////////////////////

func ExportToFile(client *elasticsearch.Client, index, filename string) error {
	scroll := time.Minute
	size := 5000

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)

	// 第一次搜索
	res, err := client.Search(
		client.Search.WithIndex(index),
		client.Search.WithSize(size),
		client.Search.WithScroll(scroll),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var sr SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return err
	}

	// 写入第一批
	for _, hit := range sr.Hits.Hits {
		b, _ := json.Marshal(hit.Source) // ⚡ 重新 Marshal 成单行 JSON
		writer.Write(b)
		writer.WriteByte('\n')
	}

	scrollID := sr.ScrollID
	if scrollID == "" {
		writer.Flush()
		return nil
	}

	// 循环 scroll
	for {
		scrollRes, err := client.Scroll(
			client.Scroll.WithScrollID(scrollID),
			client.Scroll.WithScroll(scroll),
		)
		if err != nil {
			return err
		}

		body, _ := io.ReadAll(scrollRes.Body)
		scrollRes.Body.Close()

		var sr2 SearchResponse
		if err := json.Unmarshal(body, &sr2); err != nil {
			return err
		}

		if len(sr2.Hits.Hits) == 0 { // no more docs
			break
		}

		for _, hit := range sr2.Hits.Hits {
			b, _ := json.Marshal(hit.Source) // ⚡ Marshal 成单行
			writer.Write(b)
			writer.WriteByte('\n')
		}

		scrollID = sr2.ScrollID
	}

	writer.Flush()

	// 清理 scroll
	client.ClearScroll(client.ClearScroll.WithScrollID(scrollID))

	return nil
}

//////////////////////////////////////////////////////
// 导入：从 NDJSON 文件逐行读取 → bulk 写入 ES
//////////////////////////////////////////////////////

func ImportFromFile(client *elasticsearch.Client, index, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var bulkBuf bytes.Buffer
	enc := json.NewEncoder(&bulkBuf)

	count := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// bulk metadata
		meta := map[string]any{
			"index": map[string]any{
				"_index": index,
			},
		}
		if err := enc.Encode(meta); err != nil {
			return err
		}

		// document
		if err := enc.Encode(json.RawMessage(line)); err != nil {
			return err
		}

		count++
		if count%5000 == 0 {
			if err := sendBulk(client, &bulkBuf); err != nil {
				return err
			}
		}
	}

	// 写入最后一批
	if bulkBuf.Len() > 0 {
		if err := sendBulk(client, &bulkBuf); err != nil {
			return err
		}
	}

	return nil
}

// ////////////////////////////////////////////////////
// Bulk writer
// ////////////////////////////////////////////////////
func sendBulk(client *elasticsearch.Client, buf *bytes.Buffer) error {
	res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bulk error: %s", body)
	}

	buf.Reset()
	return nil
}
