package cache

import (
	"github.com/gofier/framework/cache/driver/memory"
	"github.com/gofier/framework/cache/driver/redis"
	"github.com/gofier/framework/config"
)

var c ICache

func Init() {

}

func setStore(store string) (c ICache) {
	_conn := store
	if store == "default" {
		_conn = config.GetString("cache." + store)
		if _conn == "" {
			panic("cache connection parse error")
		}
	}

	switch _conn {
	case "memory":
		c = memory.NewMemory(
			config.GetString("cache.stores.memory.prefix"),
			config.GetUInt("cache.stores.memory.default_expiration_minute"),
			config.GetUInt("cache.stores.memory.cleanup_interval_minute"),
		)
	case "redis":
		connection := config.GetString("cache.stores.redis.connection")
		c = redis.NewRedis(
			config.GetString("database.redis."+connection+".host"),
			config.GetString("database.redis."+connection+".port"),
			config.GetString("database.redis."+connection+".password"),
			config.GetInt("database.redis."+connection+".database"),
			config.GetString("database.redis.options.prefix"),
		)
	default:
		panic("incorrect cache connection provided ")
	}

	return c
}

func Store(store string) ICache {
	return setStore(store)
}

func Cache() ICache {
	return c
}
