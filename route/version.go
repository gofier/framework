package route

import (
	"github.com/gofier/framework/request"

	"github.com/gin-gonic/gin"
)

type version struct {
	engine *request.Engine
	group  *gin.RouterGroup
	prefix string
}

func NewVersion(engine *request.Engine, prefix string) *version {
	ver := &version{
		engine: engine,
		prefix: prefix,
	}
	ver.group = ver.engine.Group(prefix)
	return ver
}

func (v *version) Auth(signKey string, relativePath string, groupFunc func(grp IGroup), handlers ...request.HandleFunc) {
	ginGroup := v.group.Group(relativePath, request.ConvertHandlers(append([]request.HandleFunc{}, handlers...))...)
	groupFunc(&group{
		engineHash:  v.engine.Hash(),
		RouterGroup: ginGroup,
	})
}

func (v *version) NoAuth(relativePath string, groupFunc func(grp IGroup), handlers ...request.HandleFunc) {
	ginGroup := v.group.Group(relativePath, request.ConvertHandlers(handlers)...)
	groupFunc(&group{
		engineHash:  v.engine.Hash(),
		RouterGroup: ginGroup,
	})
}
