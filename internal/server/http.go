package server

import (
	"demoserveice/internal/conf"
	"demoserveice/internal/service"
	"demoserveice/middleware/validate"
	"demoserveice/pkg/encode"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "gitlab.cqrb.cn/shangyou_mic/testpg/api/demoserveice/v1"
)

// NewHTTPServer new a HTTP server.

func NewHTTPServer(c *conf.Server, GreeterService *service.GreeterService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			//tracing.Server(),
			logging.Server(log.GetLogger()),
			//metrics.Server(),
			validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ErrorEncoder(encode.ErrorEncoder))
	//opts = append(opts, http.ResponseEncoder(encode.ResponseEncoder))
	srv := http.NewServer(opts...)
	r := srv.Route("")
	r.GET("/checkHealth", func(ctx http.Context) error {
		ctx.JSON(200, map[string]string{"status": "UP"})
		return nil
	})
	v1.RegisterGreeterHTTPServer(srv, GreeterService)
	return srv
}
