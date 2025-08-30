package data

import (
	"demoserveice/internal/conf"
	gormlog "demoserveice/pkg/log"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logleve "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedis, NewGRPCClient, NewHttpClient, NewMQClient, NewGreeterRepo, NewNewsRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	if c.Database.Logprefix == "" {
		return nil, nil, errors.New("请配置数据库日志文件前缀")
	}
	dsn := c.Database.Source + "?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	newLogger := gormlog.NewGormLogger("./logs/project.log", []interface{}{"service.name", c.Database.Logprefix}, logleve.Error)
	db, errDb := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`,加S
			SingularTable: true,
		},
		//对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，你可以在初始化时禁用这种方式
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if errDb != nil {
		log.Error(errDb.Error())
		return nil, nil, errDb
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(errDb.Error())
		return nil, nil, errDb
	}
	//设置连接池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	//可以复用的最大时间
	sqlDB.SetConnMaxLifetime(1700 * time.Second)
	d := &Data{
		db: db,
	}
	cleanup := func() {
		log.Info("message", "closing the data resources")
		if err := sqlDB.Close(); err != nil {
			log.Error(err)
		}
	}
	return d, cleanup, nil
}
