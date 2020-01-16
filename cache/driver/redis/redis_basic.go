package redis

import (
	"github.com/gofier/framework/cache/utils"
	"github.com/gofier/framework/zone"

	r "github.com/go-redis/redis"
)

type redisBasic struct {
	cache  *r.Client
	prefix string
}

func (re *redisBasic) Cache() *r.Client {
	return re.cache
}

func (re *redisBasic) Prefix() string {
	return re.prefix
}

func (re *redisBasic) Has(key string) bool {
	k := utils.NewKey(key, re.Prefix())

	exists, err := re.cache.Exists(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if exists <= 0 {
		return false
	}
	return true
}

func (re *redisBasic) Get(key string, defaultValue ...interface{}) interface{} {
	k := utils.NewKey(key, re.Prefix())

	if !re.Has(k.Raw()) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	valStr, err := re.cache.Get(k.Prefixed()).Result()
	if err != nil {
		return err
	}

	return valStr
}

func (re *redisBasic) Pull(key string, defaultValue ...interface{}) interface{} {
	k := utils.NewKey(key, re.Prefix())

	val := re.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}

	re.Forget(k.Raw())

	return val
}

func (re *redisBasic) Put(key string, value interface{}, future zone.Time) bool {
	k := utils.NewKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, utils.DurationFromNow(future)).Result()
	if err != nil {
		return false
	}

	return true
}

func (re *redisBasic) Add(key string, value interface{}, future zone.Time) bool {
	k := utils.NewKey(key, re.Prefix())

	// if expired return false
	ttl, err := re.cache.TTL(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if ttl > 0 {
		return false
	}

	// if exists return false
	if re.Has(k.Raw()) {
		return false
	}

	result := re.Put(k.Raw(), value, future)

	return result
}

func (re *redisBasic) Increment(key string, value int64) (incremented int64, success bool) {
	k := utils.NewKey(key, re.Prefix())

	incremented, err := re.cache.IncrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return incremented, true
}

func (re *redisBasic) Decrement(key string, value int64) (decremented int64, success bool) {
	k := utils.NewKey(key, re.Prefix())

	decremented, err := re.cache.DecrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return decremented, true
}

func (re *redisBasic) Forever(key string, value interface{}) bool {
	k := utils.NewKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, 0).Result()
	if err != nil {
		return false
	}

	return true
}

func (re *redisBasic) Forget(key string) bool {
	k := utils.NewKey(key, re.Prefix())

	result, err := re.cache.Del(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if result <= 0 {
		return false
	}

	return true
}

func (re *redisBasic) Close() error {
	return re.cache.Close()
}
