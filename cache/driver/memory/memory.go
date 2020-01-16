package memory

import (
	"errors"

	"github.com/gofier/framework/cache/utils"
	"github.com/gofier/framework/zone"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	c "github.com/patrickmn/go-cache"
)

func NewMemory(prefix string, defaultExpirationMinute uint, cleanUpIntervalMinute uint) *memory {
	return &memory{
		memoryBasic{
			cache:  c.New(zone.Duration(defaultExpirationMinute)*zone.Minute, zone.Duration(cleanUpIntervalMinute)*zone.Minute),
			prefix: prefix,
		},
	}
}

type memory struct {
	memoryBasic
}

func (m *memory) PGet(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := utils.NewKey(key, m.Prefix())

	valueInterface, found := m.cache.Get(k.Prefixed())

	if !found {
		if len(defaultValuePtr) > 0 {
			return copier.Copy(valuePtr, defaultValuePtr[0])
		}
		return errors.New("key not exist")
	}

	valueBytes, ok := valueInterface.([]byte)
	if !ok {
		return errors.New("key's value is not a valid proto buffer")
	}
	if err := proto.Unmarshal(valueBytes, valuePtr); err != nil {
		return err
	}
	return nil
}

func (m *memory) PPull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := utils.NewKey(key, m.Prefix())

	err := m.PGet(k.Raw(), valuePtr, defaultValuePtr...)
	if err != nil {
		return err
	}
	m.Forget(k.Raw())

	return nil
}

func (m *memory) PPut(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return m.Put(key, valueBytes, future)
}

func (m *memory) PAdd(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return m.Add(key, valueBytes, future)
}

func (m *memory) PForever(key string, valuePtr proto.Message) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return m.Forever(key, valueBytes)
}
