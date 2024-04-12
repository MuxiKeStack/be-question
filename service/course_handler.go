package service

import (
	"context"
	coursev1 "github.com/MuxiKeStack/be-api/gen/proto/course/v1"
	feedv1 "github.com/MuxiKeStack/be-api/gen/proto/feed/v1"
	"github.com/MuxiKeStack/be-question/domain"
	"github.com/MuxiKeStack/be-question/event"
	"strconv"
)

type CourseHandler struct {
	courseSvc coursev1.CourseServiceClient
	producer  event.Producer
}

func NewCourseHandler(courseSvc coursev1.CourseServiceClient, producer event.Producer) *CourseHandler {
	return &CourseHandler{courseSvc: courseSvc, producer: producer}
}

func (c *CourseHandler) GetRecommendationInviteeUids(ctx context.Context, question domain.Question, curUid int64, limit int64) ([]int64, error) {
	res, err := c.courseSvc.GetSubscriberUidsById(ctx, &coursev1.GetSubscriberUidsByIdRequest{
		CourseId: question.BizId,
		CurUid:   curUid,
		Limit:    limit,
	})
	if err != nil {
		return nil, err
	}
	return res.InviteeUids, nil
}

func (c *CourseHandler) InviteUserToAnswer(ctx context.Context, inviter int64, invitees []int64, question domain.Question) error {
	events := make([]event.FeedEvent, 0, len(invitees))
	for _, invitee := range invitees {
		events = append(events, event.FeedEvent{
			Type: feedv1.EventType_InviteToAnswer, // 我觉得这里可以定义枚举值
			Metadata: map[string]string{
				"inviter":    strconv.FormatInt(inviter, 10),
				"invitee":    strconv.FormatInt(invitee, 10),
				"biz":        question.Biz.String(), // 传出当前服务，则该枚举值变为易于理解的string
				"bizId":      strconv.FormatInt(question.BizId, 10),
				"questionId": strconv.FormatInt(question.Id, 10),
			},
		})
	}
	err := c.producer.BatchProduceFeedEvent(ctx, events)
	// TODO 到达feed的消费者，消费成功，要把 course 和 question 预热一下，产生一条预热的消息，然后course和question消费掉，以进行预热
	return err
}
