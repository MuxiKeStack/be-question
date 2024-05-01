package service

import (
	"context"
	"fmt"
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"github.com/MuxiKeStack/be-question/domain"
	"github.com/MuxiKeStack/be-question/repository"
)

var ErrQuestionNotFound = repository.ErrQuestionNotFound

type QuestionService interface {
	GetRecommendationInviteeUids(ctx context.Context, questionId int64, curUid int64, limit int64) ([]int64, error)
	Publish(ctx context.Context, question domain.Question) (int64, error)
	GetDetailById(ctx context.Context, questionId int64) (domain.Question, error)
	InviteUserToAnswer(ctx context.Context, inviter int64, invitees []int64, questionId int64) error
	CountBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64) (int64, error)
	ListBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64, curQuestionId int64, limit int64) ([]domain.Question, error)
	ListUserQuestions(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]domain.Question, error)
}

type questionService struct {
	repo       repository.QuestionRepository
	handlerMap map[questionv1.Biz]Handler
}

func NewQuestionService(repo repository.QuestionRepository, handlerMap map[questionv1.Biz]Handler) QuestionService {
	return &questionService{repo: repo, handlerMap: handlerMap}
}

// GetRecommendationInviteeUids 这个需要对不同的biz采取不同的方案，所有路由的一下handler
func (s *questionService) GetRecommendationInviteeUids(ctx context.Context, questionId int64, curUid int64, limit int64) ([]int64, error) {
	// get biz
	q, err := s.repo.FindById(ctx, questionId)
	if err != nil {
		return nil, err
	}
	handler, ok := s.handlerMap[q.Biz]
	if !ok {
		return nil, fmt.Errorf("未找到具体的业务处理逻辑 %s", q.Biz)
	}
	return handler.GetRecommendationInviteeUids(ctx, q, curUid, limit)
}

func (s *questionService) Publish(ctx context.Context, question domain.Question) (int64, error) {
	return s.repo.Create(ctx, question)
}

func (s *questionService) GetDetailById(ctx context.Context, questionId int64) (domain.Question, error) {
	return s.repo.FindById(ctx, questionId)
}

func (s *questionService) InviteUserToAnswer(ctx context.Context, inviter int64, invitees []int64, questionId int64) error {
	// 发消息到kafka
	q, err := s.repo.FindById(ctx, questionId)
	if err != nil {
		return err
	}
	handler, ok := s.handlerMap[q.Biz]
	if !ok {
		return fmt.Errorf("未找到具体的业务处理逻辑 %s", q.Biz)
	}
	return handler.InviteUserToAnswer(ctx, inviter, invitees, q)
}

func (s *questionService) CountBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64) (int64, error) {
	return s.repo.CountBizQuestions(ctx, biz, bizId)
}

func (s *questionService) ListBizQuestions(ctx context.Context, biz questionv1.Biz, bizId int64, curQuestionId int64, limit int64) ([]domain.Question, error) {
	return s.repo.ListBizQuestions(ctx, biz, bizId, curQuestionId, limit)
}

func (s *questionService) ListUserQuestions(ctx context.Context, uid int64, curQuestionId int64, limit int64) ([]domain.Question, error) {
	return s.repo.ListUserQuestions(ctx, uid, curQuestionId, limit)
}
