package memory

import (
	"encoding"

	"github.com/gofier/framework/cache/utils"
	"github.com/gofier/framework/zone"

	c "github.com/patrickmn/go-cache"
)

type memoryBasic struct {
	cache  *c.Cache
	prefix string
}

func (m *memoryBasic) Prefix() string {
	return m.prefix
}

func (m *memoryBasic) Has(key string) bool {
	k := utils.NewKey(key, m.Prefix())

	_, found := m.cache.Get(k.Prefixed())
	return found
}

func (m *memoryBasic) Get(key string, defaultValue ...interface{}) interface{} {
	k := utils.NewKey(key, m.prefix)

	val, found := m.cache.Get(k.Prefixed())
	if !found {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return val
}

func (m *memoryBasic) Forget(key string) bool {
	k := utils.NewKey(key, m.Prefix())

	m.cache.Delete(k.Prefixed())

	return true
}

func (m *memoryBasic) Pull(key string, defaultValue ...interface{}) interface{} {
	k := utils.NewKey(key, m.Prefix())

	val := m.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}
	m.Forget(k.Raw())
	return val
}

func (m *memoryBasic) parseValue(v interface{}) interface{} {
	switch v := v.(type) {
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil {
			return v
		}
		return b
	default:
		return v
	}
}

func (m *memoryBasic) Put(key string, value interface{}, future zone.Time) bool {
	k := utils.NewKey(key, m.Prefix())
	m.cache.Set(k.Prefixed(), m.parseValue(value), utils.DurationFromNow(future))

	return true
}

func (m *memoryBasic) Add(key string, value interface{}, future zone.Time) bool {
	k := utils.NewKey(key, m.Prefix())

	if err := m.cache.Add(k.Prefixed(), m.parseValue(value), utils.DurationFromNow(future)); err != nil {
		return false
	}

	return true
}

func (m *memoryBasic) Increment(key string, value int64) (incremented int64, success bool) {
	k := utils.NewKey(key, m.Prefix())

	incremented, err := m.cache.IncrementInt64(k.Prefixed(), m.parseValue(value).(int64))
	if err != nil {
		return 0, false
	}

	return incremented, true
}

func (m *memoryBasic) Decrement(key string, value int64) (decremented int64, success bool) {
	k := utils.NewKey(key, m.Prefix())
	decremented, err := m.cache.DecrementInt64(k.Prefixed(), m.parseValue(value).(int64))
	if err != nil {
		return 0, false
	}
	return decremented, true
}

func (m *memoryBasic) Forever(key string, value interface{}) bool {
	k := utils.NewKey(key, m.Prefix())
	m.cache.Set(k.Prefixed(), m.parseValue(value), -1)
	return true
}

func (m *memoryBasic) Close() error {
	m.cache.Flush()
	return nil
}
