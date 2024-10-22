package data

import (
	"context"
	"demoserveice/internal/conf"
	registry "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"time"
)

type HTTPClient struct {
	news *transhttp.Client
}

func NewHttpClient(c *conf.Server, logger log.Logger) *HTTPClient {
	//transhttp.WithDiscovery()
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
	conn, err := NewHttpConn(cli, "shangyou_go_news.http", "news")
	if err != nil {
		log.Error(err.Error())
	}

	d := &HTTPClient{
		news: conn,
	}
	return d
}
func NewHttpConn(cli naming_client.INamingClient, servicename string, group string) (*transhttp.Client, error) {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("discovery:///"+servicename),
		transhttp.WithDiscovery(registry.New(cli, registry.WithGroup(group))),
		transhttp.WithTimeout(500*time.Millisecond),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
