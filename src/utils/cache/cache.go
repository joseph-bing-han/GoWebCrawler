package cache

import (
	"GoWebCrawler/src/utils/conf"
	_ "fmt"
	"github.com/joseph-bing-han/gocache/cache"
	"github.com/joseph-bing-han/gocache/store"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

var (
	_cache = &cache.Cache{}
)

func init() {
	during, _ := strconv.Atoi(conf.Get("CACHE_TIME", "60000"))
	redisStore := store.NewRedis(redis.NewClient(
		&redis.Options{
			Addr:     conf.Get("REDIS_SERVER", "127.0.0.1:6379"),
			Password: conf.Get("REDIS_PASSWORD", ""),
		}),
		&store.Options{
			Expiration: time.Duration(during) * time.Minute,
		},
	)
	_cache = cache.New(redisStore)
}

func Get(key interface{}) (interface{}, error) {
	return _cache.Get(key)
}

func Has(key interface{}) bool {
	value, error := Get(key)
	if (error == nil || (error != nil && error.Error() == "redis: nil")) && value.(string) == "" {
		return false
	}
	return true
}

func Set(key interface{}, object interface{}, options ...*store.Options) error {
	if len(options) == 0 {
		return _cache.Set(key, object, nil)
	} else {
		return _cache.Set(key, object, options[0])
	}

}

func Delete(key interface{}) error {
	return _cache.Delete(key)
}

func Clear() error {
	return _cache.Clear()
}
