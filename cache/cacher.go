package cache

import "github.com/gofier/framework/cache/driver"

type ICache interface {
	driver.IProtoCache
	driver.IBasicCache
}
