package redis

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func (r *Redis) GetChannelFirstListIds(ctx context.Context, channelId int32, num int64) ([]string, error) {
	classKey := channelFirstCahceNamePre + fmt.Sprint(channelId)
	rdb := r.Source
	listId, err := rdb.ZRevRange(ctx, classKey, 0, num-1).Result()

	if err != nil {
		return nil, err
	}
	adList, err := r.GetChannelFirstAdList(ctx, channelId)
	if err != nil {
		err = errors.Wrap(err, "获取头条广告缓存失败")
		//log.Error(err)
	}
	if adList != nil {
		for _, v := range adList {
			if listId == nil {
				listId = append(listId, v.Member.(string))
			} else {
				position := int(v.Score) - 1
				//slice index 越界处理
				if position >= len(listId) {
					listId = append(listId, v.Member.(string))
					continue
				}
				rear := append([]string{}, listId[position:]...)
				listId = append(listId[0:position], v.Member.(string))
				listId = append(listId, rear...)
			}
		}
	}
	return listId, nil
}

func (r *Redis) GetChannelFirstAdList(ctx context.Context, channelId int32) ([]redis.Z, error) {
	redisKey := channelAdFirstCahceNamePre + fmt.Sprint(channelId)
	rdb := r.Source
	listId, err := rdb.ZRangeWithScores(ctx, redisKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return listId, nil
}

func (r *Redis) GetChannelTopListIds(ctx context.Context, classid int32, page int32, pagesize int32) (map[int32]string, error) {
	classKey := channelTopCahceNamePre + fmt.Sprint(classid)
	re := make(map[int32]string)
	rdb := r.Source
	offset := (page-1)*pagesize + 1
	endNum := offset + pagesize - 1

	opt := redis.ZRangeBy{Min: fmt.Sprint(offset), Max: fmt.Sprint(endNum), Count: int64(pagesize)}
	lists, err := rdb.ZRangeByScoreWithScores(ctx, classKey, &opt).Result()

	if err != nil {
		return nil, err
	}
	for _, v := range lists {
		k := int32(v.Score)
		re[k] = v.Member.(string)
	}
	return re, nil
}

func (r *Redis) GetChannelNormalListIds(ctx context.Context, classid int32, page int32, pagesize int32) ([]string, error) {
	rdb := r.Source
	classKey := channelListCahceNamePre + strconv.Itoa(int(classid))
	offset := (page - 1) * pagesize
	endNum := offset + pagesize - 1
	var listId []string
	result, err := rdb.Exists(ctx, classKey).Result()
	if result == 1 {
		listId, err = rdb.ZRevRange(ctx, classKey, int64(offset), int64(endNum)).Result()

		if err != nil {
			return nil, err
		}
	}
	return listId, nil
}
