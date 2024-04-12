package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type QuestionDAO interface {
	FindById(ctx context.Context, questionId int64) (Question, error)
	Insert(ctx context.Context, question Question) (int64, error)
}

type GORMQuestionDAO struct {
	db *gorm.DB
}

func NewGORMQuestionDAO(db *gorm.DB) QuestionDAO {
	return &GORMQuestionDAO{db: db}
}

func (dao *GORMQuestionDAO) FindById(ctx context.Context, questionId int64) (Question, error) {
	var q Question
	err := dao.db.WithContext(ctx).
		Where("id = ?", questionId).
		First(&q).Error
	return q, err
}

func (dao *GORMQuestionDAO) Insert(ctx context.Context, question Question) (int64, error) {
	now := time.Now().UnixMilli()
	question.Utime = now
	question.Ctime = now
	err := dao.db.WithContext(ctx).Create(&question).Error
	return question.Id, err
}

type Question struct {
	Id           int64 `gorm:"primaryKey,autoIncrement"`
	QuestionerId int64
	Biz          int32
	// 这个id其实应该限制住，必须是已有的资源，已有
	// 资源删除，这个连带删除的我觉得，不过因为biz的存在没办法设置外键，
	// TODO 所以这是一个要解决，但是不好解决的地方，总部查一遍再存吧，就算这样，删的时候呢，很麻烦的...
	// 可以开个定时任务吧，定期删掉错误的数据...，好像有错误数据也没啥影响...
	BizId   int64
	Content string `gorm:"varchar(200)"` // 产品给的
	Utime   int64
	Ctime   int64
}
