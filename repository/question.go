package repository

import (
	"context"
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"github.com/MuxiKeStack/be-question/domain"
	"github.com/MuxiKeStack/be-question/pkg/logger"
	"github.com/MuxiKeStack/be-question/repository/cache"
	"github.com/MuxiKeStack/be-question/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type QuestionRepository interface {
	FindById(ctx context.Context, questionId int64) (domain.Question, error)
	Create(ctx context.Context, question domain.Question) (int64, error)
	CountBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64) (int64, error)
	ListBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64, curQuestionId int64, limit int64) ([]domain.Question, error)
	ListUserQuestions(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]domain.Question, error)
}

type CachedQuestionRepository struct {
	dao   dao.QuestionDAO
	cache cache.QuestionCache
	l     logger.Logger
}

func NewCachedQuestionRepository(dao dao.QuestionDAO, cache cache.QuestionCache, l logger.Logger) QuestionRepository {
	return &CachedQuestionRepository{dao: dao, cache: cache, l: l}
}

func (repo *CachedQuestionRepository) FindById(ctx context.Context, questionId int64) (domain.Question, error) {
	res, err := repo.cache.Get(ctx, questionId)
	if err == nil {
		return res, nil
	}
	if err != cache.ErrKeyNotExist {
		// 如果数据库撑不住无缓存，全部查库，就在这里return
	}
	q, err := repo.dao.FindById(ctx, questionId)
	// 不回写，因为这是个低频操作，上面查的缓存是从发布时预热拿来的
	return repo.toDomain(q), err
}

func (repo *CachedQuestionRepository) Create(ctx context.Context, question domain.Question) (int64, error) {
	// 预缓存一下这个问题
	id, err := repo.dao.Insert(ctx, repo.toEntity(question))
	if err != nil {
		return 0, err
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := repo.cache.Set(ctx, question)
		if er != nil {
			repo.l.Error("预缓存问题失败", logger.Error(err))
		}
	}()
	return id, nil
}

func (repo *CachedQuestionRepository) CountBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64) (int64, error) {
	return repo.dao.CountBizQuestions(ctx, int32(biz), bizId)
}

func (repo *CachedQuestionRepository) ListBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64, curQuestionId int64, limit int64) ([]domain.Question, error) {
	questions, err := repo.dao.FindByBiz(ctx, int32(biz), bizId, curQuestionId, limit)
	return slice.Map(questions, func(idx int, src dao.Question) domain.Question {
		return repo.toDomain(src)
	}), err
}

func (repo *CachedQuestionRepository) ListUserQuestions(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]domain.Question, error) {
	questions, err := repo.dao.FindByUser(ctx, uid, curQuestionId, limit)
	return slice.Map(questions, func(idx int, src dao.Question) domain.Question {
		return repo.toDomain(src)
	}), err
}

func (repo *CachedQuestionRepository) toDomain(q dao.Question) domain.Question {
	return domain.Question{
		Id:           q.Id,
		QuestionerId: q.QuestionerId,
		Biz:          questionv1.Biz(q.Biz),
		BizId:        q.BizId,
		Content:      q.Content,
		Utime:        time.UnixMilli(q.Utime),
		Ctime:        time.UnixMilli(q.Ctime),
	}
}

func (repo *CachedQuestionRepository) toEntity(q domain.Question) dao.Question {
	return dao.Question{
		Id:           q.Id,
		QuestionerId: q.QuestionerId,
		Biz:          int32(q.Biz),
		BizId:        q.BizId,
		Content:      q.Content,
	}
}
