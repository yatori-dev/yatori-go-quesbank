package aiq

import (
	"encoding/json"
	"fmt"
	"yatori-go-quesbank/ques-core/entity"
	"yatori-go-quesbank/ques-core/entity/qtype"
)

// 单选题处理策略
func handleSingleChoice(topic entity.Question) AIChatMessages {
	problem := buildProblemHeader("单选题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `接下来你只需要回答选项对应内容即可...回答格式严格遵循json范式，比如：["选项1内容"]`},
		{Role: "system", Content: "就算你不知道选什么也随机选，答题过程中你无需回答任何解释，只需按照预定json范式输出！！！"},
		{Role: "system", Content: exampleSingleChoice()},
		{Role: "system", Content: problem},
	}}
}

// 构建AI题目消息
func BuildAiQuestionMessage(topic entity.Question) AIChatMessages {
	switch qtype.Index(topic.Type) {
	case qtype.SingleChoice: //单选题
		return handleSingleChoice(topic)
	case qtype.MultipleChoice: //多选题
		return handleMultipleChoice(topic)
	case qtype.TrueOrFalse: //判断题
		return handleTrueFalse(topic)
	case qtype.FillInTheBlank: //填空题
		return handleFillInTheBlank(topic)
	case qtype.ShortAnswer:
		return handleShortAnswer(topic)
	}
	return AIChatMessages{}
}

// 将AI的回复转换到Question里面
func ResponseTurnQuestion(question entity.Question, response string) entity.Question {
	var answers []string
	json.Unmarshal([]byte(response), &answers)
	question.Answers = answers
	return question
}

// 多选题处理策略
func handleMultipleChoice(topic entity.Question) AIChatMessages {
	problem := buildProblemHeader("多选题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: "接下来你只需要回答选项对应内容即可...格式：[\"选项1\",\"选项2\"]"},
		{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleMultipleChoice()},
		{Role: "system", Content: problem},
	}}
}

// 判断题处理策略
func handleTrueFalse(topic entity.Question) AIChatMessages {
	problem := buildProblemHeader("判断题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: "接下来你只需要回答“正确”或者“错误”即可...格式：[\"正确\"]"},
		{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleTrueFalse()},
		{Role: "system", Content: problem},
	}}
}

// 填空题处理策略
func handleFillInTheBlank(topic entity.Question) AIChatMessages {
	problem := buildProblemHeader("填空题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: "其中，“（answer_数字）”相关字样的地方是你需要填写答案的地方...格式：[\"答案1\",\"答案2\"]"},
		{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleFillInTheBlank()},
		{Role: "system", Content: problem},
	}}
}

// 简答题处理策略
func handleShortAnswer(topic entity.Question) AIChatMessages {
	problem := buildProblemHeader("简答题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: "这是一个简答题...格式：[\"答案\"]，注意不要拆分答案！！！"},
		{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleShortAnswer()},
		{Role: "system", Content: problem},
	}}
}

// 构建题目头部信息
func buildProblemHeader(topicType string, topic entity.Question) string {
	selectStr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	sprintf := fmt.Sprintf(`题目类型：%s
题目内容：
%s\n`, topicType, topic.Content)
	switch qtype.Index(topic.Type) {
	case qtype.SingleChoice: //单选题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + option + "\n"
		}
	case qtype.MultipleChoice: //多选题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + option + "\n"
		}
	case qtype.TrueOrFalse: //判断题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + option + "\n"
		}
	case qtype.FillInTheBlank: //填空题
	case qtype.ShortAnswer:
	}

	return sprintf
}

// 单选题示例
func exampleSingleChoice() string {
	return `比如：
试卷名称：考试
题目类型：单选
题目内容：新中国是什么时候成立的
A. 1949年10月5日
B. 1949年10月1日
C. 1949年09月1日
D. 2002年10月1日

那么你应该回答选项B的内容："["1949年10月1日"]"。注意不要携带A，B，C，D等选项前缀。`
}

// 多选题示例
func exampleMultipleChoice() string {
	return `比如：
试卷名称：考试
题目类型：多选题
题目内容：马克思关于资本积累的学说是剩余价值理论的重要组成部分...
A. 资本主义扩大再生产的源泉
B. 资本有机构成呈现不断降低趋势的根本原因
C. 社会财富占有两极分化的重要原因
D. 资本主义社会失业现象产生的根源

那么你应该回答选项A、B、D的内容："["资本主义扩大再生产的源泉","社会财富占有两极分化的重要原因","资本主义社会失业现象产生的根源"]注意不要携带A，B，C，D等选项前缀。"`
}

// 判断题示例
func exampleTrueFalse() string {
	return `比如：
试卷名称：考试
题目类型：判断
题目内容：新中国是什么时候成立是1949年10月1日吗？
A. 正确
B. 错误

那么你应该回答选项A的内容："["正确"]"。注意不要携带A，B，C，D等选项前缀。`
}

// 填空题示例
func exampleFillInTheBlank() string {
	return ` 比如：
试卷名称：考试
题目类型：填空
题目内容：新中国成立于（ ）年。
答案：1949

那么你应该回答："["1949"]"`
}

func exampleShortAnswer() string {
	return `比如：
试卷名称：考试
题目类型：简答
题目内容：请简述中国和外国的国别 differences
答案：中国和外国的国别 differences

那么你应该回答： "["中国和外国的国别 differences"]"`
}
