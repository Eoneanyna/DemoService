package data

import (
	"context"
	"demoserveice/internal/data/models/db"
	"demoserveice/internal/data/models/redis"
	"github.com/go-kratos/kratos/v2/log"
)

type NewsRepo interface {
	GetNewsById(ctx context.Context, id int32) (db.News, error)
}

type newsRepo struct {
	data  *Data
	redis *Redis
	log   *log.Helper
}

func NewNewsRepo(data *Data, redis *Redis, logger log.Logger) NewsRepo {
	return &newsRepo{
		data:  data,
		redis: redis,
		log:   log.NewHelper(logger),
	}
}

// GetNewsById 查询mysql的新闻详情
func (r *newsRepo) GetNewsById(ctx context.Context, id int32) (db.News, error) {
	dbD := db.NewDb(r.data.db)
	rdb := redis.NewRedis(r.redis.rdb)
	// 先从Redis查询
	news, err := rdb.GetNewsDetailById(ctx, int64(id))
	if err == nil {
		r.log.Infof("从Redis获取新闻详情，ID: %d", id)
		return news, nil
	}

	// Redis未命中，从MySQL查询
	r.log.Infof("Redis未命中，从MySQL获取新闻详情，ID: %d", id)
	mysqlNews, err := dbD.GetNewsById(ctx, id)
	if err != nil {
		return db.News{}, err
	}

	if len(mysqlNews) == 0 {
		return db.News{}, nil
	}

	// 将结果存入Redis缓存
	r.log.Infof("将新闻详情存入Redis缓存，ID: %d", id)
	err = rdb.SetNewsDetailById(ctx, mysqlNews[0])
	if err != nil {
		r.log.Errorf("设置新闻详情缓存失败，ID: %d, 错误: %v", id, err)
	}

	return mysqlNews[0], nil
}
