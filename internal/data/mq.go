package data

import (
	"context"
	"demoserveice/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	pb "gitlab.cqrb.cn/shangyou_mic/mq-service-pb/rocket_mq"
	"google.golang.org/grpc"
)

type MQClient struct {
	client *grpc.ClientConn
}

func NewMQClient(c *conf.Server) *MQClient {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Registry.Addr, c.Registry.Port),
	}
	cc := &constant.ClientConfig{
		NamespaceId:         c.Registry.Namespace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "nacos/cache",
		LogLevel:            c.Registry.Loglevel,
	}
	// a more graceful way to create naming client
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	conn, err := NewRpcConn(cli, "mq-service-producer.grpc", "news")
	if err != nil {
		log.Error(err.Error())
	}
	return &MQClient{client: conn}
}

type MqData struct {
	Topic string
	Msg   string
	Key   string
}

func (m *MQClient) SendMq(ctx context.Context, data *MqData) (res *pb.ProduceRepley, err error) {
	mq := pb.NewProduceClient(m.client)
	return mq.Produce(ctx, &pb.ProduceRequest{
		Topic: data.Topic,
		Msg:   data.Msg,
		Key:   data.Key,
	})
}
