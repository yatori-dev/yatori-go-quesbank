package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"yatori-go-quesbank/ques-core/entity/qtype"
)

// 用于数据库存储问题的数据结构
type DataQuestion struct {
	gorm.Model
	Question
}

// 用于回复的数据结构
type ResultQuestion struct {
	Question
	Replier string `json:"replier"` //答复者是谁
	Msg     string `json:"msg"`     //返回信息
	Code    int    `json:"code"`    //状态码，200(找到答案),404(未找到答案)
}

// 用于List的数据类型
type ListQuestion[T any] struct {
	Count   int64  `json:"count"`
	List    T      `json:"list"`
	Replier string `json:"replier"` //答复者是谁
	Msg     string `json:"msg"`     //返回信息
	Code    int    `json:"code"`    //状态码，200(找到答案),404(未找到答案)
}

type DTOQuestion struct {
	Question
	Replier   string `json:"replier"` //答复者是谁
	ReplyType string `json:"replyType"`
}

// 问题数据结构
type Question struct {
	Md5     string      `gorm:"column:md5" json:"md5"`                            //题目MD5值，注意，是（题目类型+题目内容）的编码的MD5值
	Type    qtype.QType `gorm:"column:type" json:"type" binding:"required"`       //题目类型
	Content string      `gorm:"column:content" json:"content" binding:"required"` //题目内容
	Options StringArray `gorm:"column:options" json:"options"`                    //选项（一般选择题才会有），存储为Json
	Answers StringArray `gorm:"column:answers;type:TEXT" json:"answers"`          // 答案，存储为 JSON

}

type StringArray []string

// 字符串转StringArray
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// StringArray转字符串
func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte: %T", value)
	}
	return json.Unmarshal(bytes, s)
}
