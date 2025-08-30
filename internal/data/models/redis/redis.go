package redis

import "github.com/redis/go-redis/v9"

type Redis struct {
	Source *redis.Client
}

const (
	// 频道列表缓存名前缀
	channelListCahceNamePre = "list:channelid:"
	// 频道列表缓存名前缀
	channelFirstCahceNamePre = "first:channelid:"
	// 频道置顶缓存名前缀
	channelTopCahceNamePre = "top:channelid:"
	// 频道广告缓存名前缀
	channelAdFirstCahceNamePre = "adfirst:channelid:"

	// 新闻热点列表缓存名前缀
	newsListCacheId = "list:hotSort"
	// 新闻详情的一个hash列表
	newsDetailCacheNamePre = "news:detail:"
)
