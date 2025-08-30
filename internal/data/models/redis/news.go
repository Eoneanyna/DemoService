package redis

import (
	"context"
	"demoserveice/internal/data/models/db"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// GetNewsDetailById 获取新闻详情
func (r *Redis) GetNewsDetailById(ctx context.Context, id int64) (db.News, error) {
	rdb := r.Source
	news := db.News{}

	raw, err := rdb.HGet(ctx, newsDetailCacheNamePre, fmt.Sprintf("%d", id)).Result()
	if err != nil {
		return news, err
	}

	if err = json.Unmarshal([]byte(raw), &news); err != nil {
		return news, err
	}
	return news, nil
}

// SetNewsDetailById 刷新新闻详情
func (r *Redis) SetNewsDetailById(ctx context.Context, news db.News) error {
	rdb := r.Source
	data, err := json.Marshal(news)
	if err != nil {
		return err
	}

	_, err = rdb.HSet(ctx, fmt.Sprintf("%s%d", newsDetailCacheNamePre, news.Id, data)).Result()
	if err != nil {
		return err
	}

	// TODO 设置过期时间
	rdb.Expire(ctx, newsDetailCacheNamePre, 24*time.Hour)

	return nil
}

// 获取热点排行
func (r *Redis) GetNewsHotList(ctx context.Context, key string, page int32, pagesSize int32) (map[int32]string, error) {
	rdb := r.Source

	re := make(map[int32]string)
	offset := (page-1)*pagesSize + 1
	endNum := offset + pagesSize - 1

	opt := redis.ZRangeBy{Min: fmt.Sprint(offset), Max: fmt.Sprint(endNum), Count: int64(pagesSize)}
	lists, err := rdb.ZRangeByScoreWithScores(ctx, newsListCacheId, &opt).Result()

	if err != nil {
		return nil, err
	}
	for _, v := range lists {
		k := int32(v.Score)
		re[k] = v.Member.(string)
	}
	return re, nil
}

// 刷新热点排行
func (r *Redis) SetNewsHotList(ctx context.Context, key string, newsList []*db.News, expiration time.Duration) error {
	rdb := r.Source

	// 1. 存储新闻详情
	for _, news := range newsList {
		detailKey := fmt.Sprintf("%s%d", newsListCacheId, news.Id)
		data, _ := json.Marshal(news)
		rdb.Set(ctx, detailKey, data, 24*time.Hour)

		// 根据排序类型设置不同的分数
		var score float64
		if strings.Contains(key, ":hot:") {
			score = float64(news.ViewCount) // 按点击量排序
		} else {
			score = float64(news.CreateTime) // 按时间排序
		}

		rdb.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: news.Id,
		})
	}

	// 设置过期时间
	rdb.Expire(ctx, key, expiration)
	return nil
}
