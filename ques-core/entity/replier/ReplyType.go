package replier

type ReplyType int

const (
	LOCAL ReplyType = iota
	TONGYI
	OPENAI
	DEEPSEEK
	DOUBAO
	YANXI
)

var ReplyTypeStr = [...]string{
	"本地缓存数据库",
	"通义千问",
	"ChatGPT",
	"DeepSeek",
	"豆包",
	"言溪题库",
}

// 转QType转Str
func (q ReplyType) String() string {
	return ReplyTypeStr[q]
}

// Str转QType
func Index(replyStr string) ReplyType {
	for k, v := range ReplyTypeStr {
		if v == replyStr {
			return ReplyType(k)
		}
	}
	return -1
}
