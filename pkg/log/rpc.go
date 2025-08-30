package log

import (
	"context"
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

type RegisterConf struct {
	Addr      string
	Port      uint64
	Namespace string
	Loglevel  string
}

func NewMqClient(c *RegisterConf) (*grpc.ClientConn, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Addr, c.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         c.Namespace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "nacos/cache",
		LogLevel:            c.Loglevel,
	}

	// a more graceful way to create naming client
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	mqConn, err := NewRpcConn(cli, "mq-service-producer.grpc", "news")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return mqConn, nil
}

func NewRpcConn(cli naming_client.INamingClient, servicename string, group string) (*grpc.ClientConn, error) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
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
