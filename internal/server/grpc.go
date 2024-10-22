package server

import (
	"context"
	"demoserveice/internal/conf"
	"demoserveice/internal/service"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pkg/errors"
	v12 "gitlab.cqrb.cn/shangyou_mic/testpg/api/demoserveice/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, getter *service.GreeterService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(log.GetLogger()),
			validate.Validator(),
		),
	}

	ctx := context.Background()
	err := initTracer(ctx)
	if err != nil {
		log.Error("" + err.Error())
	}
	opts = append(opts, grpc.Middleware(
		tracing.Server(),
	))
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
	return srv
}

func initTracer(ctx context.Context) error {
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("172.29.122.250:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return errors.Wrap(err, "172.16.18.214测试服Jaeger初始化失败")
	}
	tp := trace.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		trace.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("demoserveice"),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 239525.56),
		)),
	)
	otel.SetTracerProvider(tp)
	fmt.Println("initTracer ok")
	return nil
}
