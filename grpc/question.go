package grpc

import (
	"context"
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"github.com/MuxiKeStack/be-question/domain"
	"github.com/MuxiKeStack/be-question/service"
	"google.golang.org/grpc"
)

type QuestionServiceServer struct {
	questionv1.UnimplementedQuestionServiceServer
	svc service.QuestionService
}

func NewQuestionServiceServer(svc service.QuestionService) *QuestionServiceServer {
	return &QuestionServiceServer{svc: svc}
}

func (s *QuestionServiceServer) Register(server grpc.ServiceRegistrar) {
	questionv1.RegisterQuestionServiceServer(server, s)
}

func (s *QuestionServiceServer) Publish(ctx context.Context, request *questionv1.PublishRequest) (*questionv1.PublishResponse, error) {
	qid, err := s.svc.Publish(ctx, convertToDomain(request.GetQuestion()))
	return &questionv1.PublishResponse{
		QuestionId: qid,
	}, err
}

func (s *QuestionServiceServer) GetRecommendationInviteeUids(ctx context.Context, request *questionv1.GetRecommendationInviteeUidsRequest) (*questionv1.GetRecommendationInviteeUidsResponse, error) {
	uids, err := s.svc.GetRecommendationInviteeUids(ctx, request.GetQuestionId(), request.GetCurUid(), request.GetLimit())
	return &questionv1.GetRecommendationInviteeUidsResponse{InviteeUids: uids}, err
}

func (s *QuestionServiceServer) GetDetailById(ctx context.Context, request *questionv1.GetDetailByIdRequest) (*questionv1.GetDetailByIdResponse, error) {
	q, err := s.svc.GetDetailById(ctx, request.GetQuestionId())
	return &questionv1.GetDetailByIdResponse{Question: convertToVo(q)}, err
}

func (s *QuestionServiceServer) InviteUserToAnswer(ctx context.Context, request *questionv1.InviteUserToAnswerRequest) (*questionv1.InviteUserToAnswerResponse, error) {
	err := s.svc.InviteUserToAnswer(ctx, request.GetInviter(), request.GetInvitees(), request.GetQuestionId())
	return &questionv1.InviteUserToAnswerResponse{}, err
}

func convertToDomain(q *questionv1.Question) domain.Question {
	return domain.Question{
		Id:           q.GetId(),
		QuestionerId: q.GetQuestionerId(),
		Biz:          q.GetBiz(),
		BizId:        q.GetBizId(),
		Content:      q.GetContent(),
	}
}

func convertToVo(q domain.Question) *questionv1.Question {
	return &questionv1.Question{
		Id:           q.Id,
		QuestionerId: q.QuestionerId,
		Biz:          q.Biz,
		BizId:        q.BizId,
		Content:      q.Content,
		Utime:        q.Utime.UnixMilli(),
		Ctime:        q.Ctime.UnixMilli(),
	}
}
