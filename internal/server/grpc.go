package server

import (
	v1 "demoserveice/api/news/v1"
	"demoserveice/internal/conf"
	"demoserveice/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v12 "gitlab.cqrb.cn/shangyou_mic/testpg/api/demoserveice/v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, news *service.NewsService, getter *service.GreeterService, ds *service.DataStreamService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(log.GetLogger()),
			validate.Validator(),
			tracing.Server(),
		),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v12.RegisterGreeterServer(srv, getter)
	v12.RegisterDataStreamServer(srv, ds)
	v1.RegisterNewsServiceServer(srv, news)
	return srv
}
