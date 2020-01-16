package driver

import (
	"github.com/gofier/framework/zone"

	"github.com/golang/protobuf/proto"
)

type ProtoCacheGetter interface {
	PGet(key string, valuePrt proto.Message, defaultValuePtr ...proto.Message) error
}

type IProtoCache interface {
	PPull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error
	PPut(key string, valuePtr proto.Message, future zone.Time) bool
	PAdd(key string, valuePtr proto.Message, future zone.Time) bool
	PForever(key string, valuePtr proto.Message) bool

	ProtoCacheGetter
}
