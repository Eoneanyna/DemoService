package data

import (
	"context"
	"demoserveice/internal/conf"
	registry "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	news *grpc.ClientConn
}

func NewGRPCClient(c *conf.Server, logger log.Logger) *GRPCClient {
	log := log.NewHelper(logger)
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Registry.Addr, c.Registry.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         c.Registry.Namespace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: c.Registry.Loglevel,
	}

	// a more graceful way to create naming client
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	conn, err := NewRpcConn(cli, "demoservice", "news")
	if err != nil {
		log.Error(err.Error())
	}
	d := &GRPCClient{
		news: conn,
	}

	return d
}
func NewRpcConn(cli naming_client.INamingClient, servicename string, group string) (*grpc.ClientConn, error) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		//transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			middleware.Chain(
				recovery.Recovery(),
				mmd.Client(),
			),
		),
		transgrpc.WithEndpoint("discovery:///"+servicename),
		transgrpc.WithDiscovery(registry.New(cli, registry.WithGroup(group))),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
