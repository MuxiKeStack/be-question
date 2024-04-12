package service

import (
	"context"
	"github.com/MuxiKeStack/be-question/domain"
)

type Handler interface {
	GetRecommendationInviteeUids(ctx context.Context, question domain.Question, curUid int64, limit int64) ([]int64, error)
	InviteUserToAnswer(ctx context.Context, inviter int64, invitees []int64, question domain.Question) error
}
