package examples

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Article struct {
	ID      int
	Content string
}

func Test_fts(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// 创建 FTS5 虚拟表
	db.Exec("CREATE VIRTUAL TABLE articles USING fts5(content, tokenize = 'unicode61');")

	// 插入数据
	db.Exec("INSERT INTO articles (content) VALUES (?)", "Go语言教程")
	db.Exec("INSERT INTO articles (content) VALUES (?)", "Python语言入门")
	db.Exec("INSERT INTO articles (content) VALUES (?)", "Go web 开发")

	// 搜索
	var results []Article
	db.Raw(`SELECT rowid as id, content FROM articles WHERE articles MATCH ?`, "Goaf").Scan(&results)

	for _, r := range results {
		println(r.Content)
	}
}
