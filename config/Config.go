package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"yatori-go-quesbank/ques-core/entity/aitype"
)

type JSONDataForConfig struct {
	Setting Setting `json:"setting"`
}
type BasicSetting struct {
	LogOutFileSw int    `json:"logOutFileSw,omitempty" yaml:"logOutFileSw"` //是否输出日志文件0代表不输出，1代表输出，默认为1
	LogLevel     string `json:"logLevel,omitempty" yaml:"logLevel"`         //日志等级，默认INFO，DEBUG为找BUG调式用的，日志内容较详细，默认为INFO
}
type AiSetting struct {
	AiType  aitype.AiType `json:"aiType" yaml:"aiType"`
	AiUrl   string        `json:"aiUrl" yaml:"aiUrl"`
	AiModel string        `json:"aiModel"`
	APIKEY  string        `json:"API_KEY" yaml:"API_KEY" mapstructure:"API_KEY"`
}

type ExternalSetting struct {
	ExType  string `json:"exType" yaml:"exType"`   //外部题库对接类型
	ExToken string `json:"exToken" yaml:"exToken"` //外部题库对接Token
}

type AnswerSetting struct {
	AnswerType string `json:"answerType"`
	AiSetting
	ExternalSetting
}
type Setting struct {
	BasicSetting  BasicSetting    `json:"basicSetting" yaml:"basicSetting"`
	AnswerSetting []AnswerSetting `json:"answerSetting" yaml:"answerSetting"`
}

// 读取json配置文件
func ReadJsonConfig(filePath string) JSONDataForConfig {
	var configJson JSONDataForConfig
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &configJson)
	if err != nil {
		log.Fatal(err)
	}
	return configJson
}

// 自动识别读取配置文件
func ReadConfig(filePath string) JSONDataForConfig {
	var configJson JSONDataForConfig
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		//log2.Print(log2.INFO, log2.BoldRed, "找不到配置文件或配置文件内容书写错误")
		log.Fatal(err)
	}
	err = viper.Unmarshal(&configJson)
	//viper.SetTypeByDefaultValue(true)
	viper.SetDefault("setting.basicSetting.logModel", 5)

	if err != nil {
		//log2.Print(log2.INFO, log2.BoldRed, "配置文件读取失败，请检查配置文件填写是否正确")
		log.Fatal(err)
	}
	return configJson
}

// CmpCourse 比较是否存在对应课程,匹配上了则true，没有匹配上则是false
func CmpCourse(course string, courseList []string) bool {
	for i := range courseList {
		if courseList[i] == course {
			return true
		}
	}
	return false
}

func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func StrToInt(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0 // 其他错误处理逻辑
	}
	return res
}
