package ioc

import (
	coursev1 "github.com/MuxiKeStack/be-api/gen/proto/course/v1"
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"github.com/MuxiKeStack/be-question/event"
	"github.com/MuxiKeStack/be-question/service"
)

func InitSvcHandlers(courseSvc coursev1.CourseServiceClient, producer event.Producer) map[questionv1.Biz]service.Handler {
	courseHandler := service.NewCourseHandler(courseSvc, producer)
	return map[questionv1.Biz]service.Handler{
		questionv1.Biz_Course: courseHandler,
	}
}
