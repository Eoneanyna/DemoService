package main

import (
	"demoserveice/internal/conf"
	"flag"
	"fmt"
	config "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	registry "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	kitlog "gitlab.cqrb.cn/shangyou_mic/kit/log"
	"gitlab.cqrb.cn/shangyou_mic/kit/reporter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "demoserveice"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, reg *conf.Server) *kratos.App {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(reg.Registry.Addr, reg.Registry.Port),
	}

	cc := constant.ClientConfig{
		NamespaceId:         reg.Registry.Namespace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "nacos/cache",
		LogLevel:            reg.Registry.Loglevel,
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(registry.New(client, registry.WithCluster("DEFAULT"), registry.WithGroup("news"))),
	)
}

func main() {
	flag.Parse()
	configLogLevel := os.Getenv("CONFIGLOGLEVEL")
	if configLogLevel == "" {
		configLogLevel = "error"
	}
	encoder := zapcore.EncoderConfig{
		TimeKey:   "t",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	logger := kitlog.NewZapLogger(
		encoder,
		kitlog.LogLevelStr(configLogLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	logAddKey := []interface{}{
		"service.name", Name,
	}
	log.SetLogger(log.With(logger, logAddKey...))

	//读取配置文件
	confighost := os.Getenv("CONFIGHOST")
	if confighost == "" {
		panic("请先设置配置中心地址")
	}
	Namespace := os.Getenv("NAMESPACE")
	if Namespace == "" {
		Namespace = "develop"
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(confighost, 8848),
	}
	cc := constant.ClientConfig{
		NamespaceId:         Namespace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "nacos/cache",
		LogLevel:            configLogLevel,
	}
	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	c := kconfig.New(
		kconfig.WithSource(
			config.NewConfigSource(client, config.WithGroup(Name), config.WithDataID("config.yaml")),
		),
		kconfig.WithDecoder(func(kv *kconfig.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	//钉钉报错配置
	if bc.Reporter.IsReport && bc.Reporter.Systemid != 0 {
		rpo := reporter.NewRepoter(bc.Reporter.Systemid, bc.Reporter.Endpoint)
		logger.WithReporter(rpo)
	}
	//日志收集服务初始化 namespace 为mq所在的命名空间(develop|testing|production), 与当前命名空间不一致时，需要修改

	if err != nil {
		panic(err)
	}

	fmt.Println("http端口", bc.Server.Http.Addr)
	fmt.Println("grpc端口", bc.Server.Grpc.Addr)
	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
