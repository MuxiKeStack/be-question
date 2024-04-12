//go:build wireinject

package main

import (
	"github.com/MuxiKeStack/be-question/grpc"
	"github.com/MuxiKeStack/be-question/ioc"
	"github.com/MuxiKeStack/be-question/pkg/grpcx"
	"github.com/MuxiKeStack/be-question/repository"
	"github.com/MuxiKeStack/be-question/repository/cache"
	"github.com/MuxiKeStack/be-question/repository/dao"
	"github.com/MuxiKeStack/be-question/service"
	"github.com/google/wire"
)

func InitGRPCServer() grpcx.Server {
	wire.Build(
		ioc.InitGRPCxKratosServer,
		grpc.NewQuestionServiceServer,
		service.NewQuestionService,
		ioc.InitSvcHandlers,
		ioc.InitCourseClient,
		ioc.InitProducer,
		repository.NewCachedQuestionRepository,
		dao.NewGORMQuestionDAO, cache.NewRedisQuestionCache,
		// 第三方
		ioc.InitKafka,
		ioc.InitDB,
		ioc.InitEtcdClient,
		ioc.InitLogger,
		ioc.InitRedis,
	)
	return grpcx.Server(nil)
}
