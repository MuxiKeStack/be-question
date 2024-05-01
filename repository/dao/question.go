package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type QuestionDAO interface {
	FindById(ctx context.Context, questionId int64) (Question, error)
	Insert(ctx context.Context, question Question) (int64, error)
	CountBizQuestions(ctx context.Context, biz int32, bizId int64) (int64, error)
	FindByBiz(ctx context.Context, biz int32, bizId int64, curQuestionId int64, limit int64) ([]Question, error)
	FindByUser(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]Question, error)
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

func (dao *GORMQuestionDAO) CountBizQuestions(ctx context.Context, biz int32, bizId int64) (int64, error) {
	var cnt int64
	err := dao.db.WithContext(ctx).
		Model(&Question{}).
		Where("biz = ? and biz_id = ?", biz, bizId).
		Count(&cnt).Error
	return cnt, err
}

func (dao *GORMQuestionDAO) FindByBiz(ctx context.Context, biz int32, bizId int64, curQuestionId int64, limit int64) ([]Question, error) {
	var qs []Question
	err := dao.db.WithContext(ctx).
		// question是不可更改的，所以新旧排序可以按照id来排序
		Where("biz = ? and biz_id = ? and id < ?", biz, bizId, curQuestionId).
		Order("id desc").
		Limit(int(limit)).
		Find(&qs).Error
	return qs, err
}

func (dao *GORMQuestionDAO) FindByUser(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]Question, error) {
	var qs []Question
	err := dao.db.WithContext(ctx).
		// question是不可更改的，所以新旧排序可以按照id来排序
		Where("questioner_id = ? and id < ?", uid, curQuestionId).
		Order("id desc").
		Limit(int(limit)).
		Find(&qs).Error
	return qs, err
}

type Question struct {
	Id           int64 `gorm:"primaryKey,autoIncrement"`
	QuestionerId int64 `gorm:"index"`
	Biz          int32 `gorm:"index:biz_bizId"`
	// 这个id其实应该限制住，必须是已有的资源，已有
	// 资源删除，这个连带删除的我觉得，不过因为biz的存在没办法设置外键，
	// TODO 所以这是一个要解决，但是不好解决的地方，总部查一遍再存吧，就算这样，删的时候呢，很麻烦的...
	// 可以开个定时任务吧，定期删掉错误的数据...，好像有错误数据也没啥影响...
	BizId   int64  `gorm:"index:biz_bizId"`
	Content string `gorm:"varchar(200)"` // 产品给的长度
	Utime   int64
	Ctime   int64
}
