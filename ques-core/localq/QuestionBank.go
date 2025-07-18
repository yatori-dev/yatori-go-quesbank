package questionbank

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"time"
	"yatori-go-quesbank/ques-core/entity"
)

// 题库缓存初始化
func QuestionBackInit() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("questionbank.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.Question{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB() //数据库连接池
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

// 插入题库
func Insert(db *gorm.DB, question *entity.DataQuestion) error {
	if err := db.Create(&question).Error; err != nil {
		return errors.New("插入数据失败: " + err.Error())
	}
	//log2.Print(log2.DEBUG, "插入数据成功")
	return nil
}

// 如果没有则插入题库
func InsertIfNot(db *gorm.DB, question *entity.DataQuestion) error {
	//检查是否合法题目
	checkErr := CheckQue(question)
	if checkErr != nil {
		return checkErr
	}
	selectQs := SelectsForTypeAndContent(db, question)
	if len(selectQs) > 0 {
		return nil
	}
	if question.Md5 == "" {
		question.Md5 = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s", question.Type, question.Content))))
	}
	// 插入题目
	err := Insert(db, question)
	if err != nil {
		return err
	}
	return nil
}

// 根据题目类型和内容查询题目
func SelectsForTypeAndContent(db *gorm.DB, question *entity.DataQuestion) []entity.DataQuestion {
	var questions []entity.DataQuestion
	if err := db.Where("type = ? AND content = ?", question.Type, question.Content).Find(&questions).Error; err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	return questions
}

// 根据题目类型和内容查询题目
func SelectForTypeAndContent(db *gorm.DB, question *entity.DataQuestion) *entity.DataQuestion {
	var qu entity.DataQuestion
	if err := db.First(&qu, "type = ? AND content = ?", question.Type, question.Content).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.Fatalf("查询数据失败: %v", err)
	}
	return &qu
}

// 根据题目类型和内容查询题目
func SelectForTypeAndLikeContent1_4(db *gorm.DB, question *entity.DataQuestion) *entity.DataQuestion {
	var qu entity.DataQuestion
	if len(question.Content) < 10 {
		if err := db.First(&qu, "type = ? AND content = ?", question.Type, question.Content).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			log.Fatalf("查询数据失败: %v", err)
		}
		return &qu
	} else {
		runes := []rune(question.Content)
		slog.Debug(string(runes[2 : len(runes)-5]))

		if err := db.Where("type = ? AND content like ?", question.Type, "%"+string(runes[2:len(runes)-5])+"%").First(&qu).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			log.Fatalf("查询数据失败: %v", err)
		}
		return &qu
	}

}

// 根据题目MD5查询
func SelectsForMd5(db *gorm.DB, question *entity.DataQuestion) []entity.DataQuestion {
	var questions []entity.DataQuestion
	if err := db.Where("md5 = ?", question.Md5).Find(&questions).Error; err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	return questions
}

// 直接通过题目找答案返回
func SelectAnswer(db *gorm.DB, question *entity.DataQuestion) []string {

	return nil
}

// 根据题目类型和内容更新题目
func UpdateAnswerForTypeAndContent(db *gorm.DB, question *entity.DataQuestion) error {
	if err := db.Where("type = ? AND content = ?", question.Type, question.Content).Updates(&question).Error; err != nil {
		return err
	}
	return nil
}

// 根据题目类型和内容删除题目
func DeleteForTypeAndContent(db *gorm.DB, question *entity.DataQuestion) error {
	if err := db.Where("type = ? AND content = ?", question.Type, question.Content).Delete(&entity.DataQuestion{}).Error; err != nil {
		return err
	}
	return nil
}

// 检验Question合法性
func CheckQue(question *entity.DataQuestion) error {
	//检验数据合法性
	if question.Type == "" {
		return errors.New("Not Found Question Type")
	}
	if question.Content == "" {
		return errors.New("Not Found Question Content")
	}
	return nil
}
