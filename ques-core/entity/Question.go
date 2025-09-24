package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// 用于数据库存储问题的数据结构
type DataQuestion struct {
	gorm.Model
	Question
	RightStatus int `gorm:"not null;column:right_status;type:INTEGER;default:0" json:"right_status"` //正确状态，0为待定，1为错误，2为正确
}

// Es用的
type EsQuestion struct {
	Question
	RightStatus int `gorm:"not null;column:right_status;type:INTEGER;default:0" json:"right_status"` //正确状态，0为待定，1为错误，2为正确
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
	Md5     string      `gorm:"column:md5" json:"md5"`                   //题目MD5值，注意，是（题目类型+题目内容）的编码的MD5值
	Type    string      `gorm:"not null;column:type" json:"type"`        //题目类型
	Content string      `gorm:"not null;column:content" json:"content"`  //题目内容
	Options StringArray `gorm:"column:options" json:"options"`           //选项（一般选择题才会有），存储为Json
	Answers StringArray `gorm:"column:answers;type:TEXT" json:"answers"` // 答案，存储为 JSON

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
func (s *Question) String() string {
	marshal, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%x", s)
	}
	return string(marshal)
}

func (s *DTOQuestion) String() string {
	marshal, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%x", s)
	}
	return string(marshal)
}
