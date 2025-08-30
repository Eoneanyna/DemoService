module demoserveice

go 1.22

replace gitlab.cqrb.cn/shangyou_mic/testpg => D:\golang\code\testpg

//replace gitlab.cqrb.cn/shangyou_mic/kit => /Users/lhb/www/goweb/src/cqrb/kit

require (
	github.com/go-kratos/kratos/contrib/config/nacos/v2 v2.0.0-20240909031108-908e6256a9f6
	github.com/go-kratos/kratos/contrib/registry/nacos/v2 v2.0.0-20240909031108-908e6256a9f6
	github.com/go-kratos/kratos/v2 v2.8.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.4
	github.com/google/wire v0.6.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/nacos-group/nacos-sdk-go v1.1.4
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.6.1
	gitlab.cqrb.cn/shangyou_mic/testpg v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.47 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/strftime v1.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.20.1 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
