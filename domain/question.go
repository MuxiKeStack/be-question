package domain

import (
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"time"
)

type Question struct {
	Id           int64
	QuestionerId int64
	Biz          questionv1.Biz
	BizId        int64
	Content      string
	Utime        time.Time
	Ctime        time.Time
}
