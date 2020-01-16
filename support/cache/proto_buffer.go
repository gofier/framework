package cache

import (
	c "github.com/gofier/framework/cache"
	"github.com/gofier/framework/zone"

	"github.com/golang/protobuf/proto"
)

func PGet(key string, valuePtr proto.Message, defaultValue ...proto.Message) error {
	return c.Cache().PGet(key, valuePtr, defaultValue...)
}
func PPull(key string, valuePtr proto.Message, defaultValue ...proto.Message) error {
	return c.Cache().PPull(key, valuePtr, defaultValue...)
}
func PPut(key string, value proto.Message, future zone.Time) bool {
	return c.Cache().PPut(key, value, future)
}
func PAdd(key string, value proto.Message, future zone.Time) bool {
	return c.Cache().PAdd(key, value, future)
}
func PForever(key string, value proto.Message) bool {
	return c.Cache().PForever(key, value)
}
