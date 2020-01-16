package redis

import (
	"errors"

	"github.com/gofier/framework/cache/utils"
	"github.com/gofier/framework/zone"

	r "github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

func NewRedis(host, port, password string, dbIndex int, prefix string) *redis {
	client := r.NewClient(&r.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbIndex,
	})

	return &redis{
		redisBasic{
			cache:  client,
			prefix: prefix,
		},
	}
}

type redis struct {
	redisBasic
}

func (re *redis) PGet(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	var valueBytes []byte

	k := utils.NewKey(key, re.Prefix())

	if !re.Has(k.Raw()) {
		if len(defaultValuePtr) > 0 {
			return copier.Copy(valuePtr, defaultValuePtr[0])
		}
		return errors.New("key not exist")
	}
	err := re.cache.Get(k.Prefixed()).Scan(&valueBytes)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(valueBytes, valuePtr); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------------------------
// the same
// ------------------------------------------------------------------------------

func (re *redis) PPull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := utils.NewKey(key, re.Prefix())

	err := re.PGet(k.Raw(), valuePtr, defaultValuePtr...)
	if err != nil {
		return err
	}

	re.Forget(k.Raw())

	return nil
}
func (re *redis) PPut(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Put(key, valueBytes, future)
}
func (re *redis) PAdd(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Add(key, valueBytes, future)
}
func (re *redis) PForever(key string, valuePtr proto.Message) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return re.Forever(key, valueBytes)
}
