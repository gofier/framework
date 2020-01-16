package route

import (
	"github.com/gofier/framework/policy"
	"github.com/gofier/framework/request"
	"github.com/gofier/framework/request/websocket"

	"github.com/gin-gonic/gin"
)

type IRouteGroup interface {
	Group(grp IGroup)
}

type IGroup interface {
	AddGroup(relativePath string, routeGrouper IRouteGroup, handlers ...request.HandleFunc)
	iRoutes
}

type iRoutes interface {
	Handle(httpMethod, relativePath string, handlers ...request.HandleFunc) routeEnder
	Any(relativePath string, handlers ...request.HandleFunc) routeEnder
	GET(relativePath string, handlers ...request.HandleFunc) routeEnder
	POST(relativePath string, handlers ...request.HandleFunc) routeEnder
	DELETE(relativePath string, handlers ...request.HandleFunc) routeEnder
	PATCH(relativePath string, handlers ...request.HandleFunc) routeEnder
	PUT(relativePath string, handlers ...request.HandleFunc) routeEnder
	OPTIONS(relativePath string, handlers ...request.HandleFunc) routeEnder
	HEAD(relativePath string, handlers ...request.HandleFunc) routeEnder

	Websocket(relativePath string, wsHandler websocket.Handler, handlers ...request.HandleFunc) routeEnder

	StaticFile(relativePath, filepath string) gin.IRoutes
	Static(relativePath, root string) gin.IRoutes
	// StaticFS(relativePath string)
}

type group struct {
	engineHash request.EngineHash
	*gin.RouterGroup
}

type routeEnder interface {
	policy.RoutePolicier
	Name(routeName string) policy.RoutePolicier
}

func (g *group) clearPath(relativePath string) string {
	if relativePath == "" {
		return relativePath
	}
	basePath := g.RouterGroup.BasePath()
	if basePath[len(basePath)-1:] == "/" && relativePath[:1] == "/" {
		return relativePath[1:]
	}
	if basePath[len(basePath)-1:] != "/" && relativePath[:1] != "/" {
		return "/" + relativePath
	}
	return relativePath
}

func (g *group) Handle(httpMethod, relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute(httpMethod, g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.Handle(httpMethod, relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) Any(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("Any", g, g.clearPath(relativePath), func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.Any(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) GET(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("GET", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.GET(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) POST(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("POST", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.POST(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) DELETE(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("DELETE", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.DELETE(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) PATCH(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("PATCH", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.PATCH(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) PUT(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("PUT", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.PUT(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) OPTIONS(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("OPTIONS", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.OPTIONS(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) HEAD(relativePath string, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("HEAD", g, relativePath, func(innerHandlers ...request.HandleFunc) {
		g.RouterGroup.HEAD(relativePath, request.ConvertHandlers(innerHandlers)...)
	}, handlers...)
}

func (g *group) StaticFile(relativePath, filepath string) gin.IRoutes {
	relativePath = g.clearPath(relativePath)
	return g.RouterGroup.StaticFile(relativePath, filepath)
}

func (g *group) Static(relativePath, root string) gin.IRoutes {
	relativePath = g.clearPath(relativePath)
	return g.RouterGroup.Static(relativePath, root)
}

// func (g *group) StaticFS(relativePath string) {
//
// }

const httpMethodWebsocket = "WS"

func (g *group) Websocket(relativePath string, wsHandler websocket.Handler, handlers ...request.HandleFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newWsRoute(httpMethodWebsocket, g, relativePath, func(wsHandler websocket.Handler, innerHandlers ...request.HandleFunc) {
		innerGinHandlers := append(request.ConvertHandlers(innerHandlers), websocket.CovertHandler(wsHandler))
		g.RouterGroup.GET(relativePath, innerGinHandlers...)
	}, wsHandler, handlers...)
}

func (g *group) AddGroup(relativePath string, routeGrouper IRouteGroup, handlers ...request.HandleFunc) {
	ginGroup := g.RouterGroup.Group(relativePath, request.ConvertHandlers(handlers)...)
	routeGrouper.Group(&group{
		engineHash:  g.engineHash,
		RouterGroup: ginGroup,
	})
}
