package route

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/gofier/framework/config"
	"github.com/gofier/framework/http/middleware"
	"github.com/gofier/framework/log"
	"github.com/gofier/framework/policy"
	"github.com/gofier/framework/request"
	"github.com/gofier/framework/request/websocket"
)

var RouteNameMap *routeNameMap

type routeNameMap struct {
	lock sync.RWMutex
	data map[string]string
}

func newRouteNameMap() *routeNameMap {
	return &routeNameMap{
		data: make(map[string]string),
	}
}

func (rnm *routeNameMap) set(routeName, routeUrl string) {
	rnm.lock.Lock()
	defer rnm.lock.Unlock()
	rnm.data[routeName] = routeUrl
}

func (rnm *routeNameMap) Get(routeName string, param map[string]string) (url string, err error) {
	rnm.lock.RLock()
	defer rnm.lock.RUnlock()
	routeUrl, ok := rnm.data[routeName]
	if !strings.Contains(routeUrl, "/:") {
		return routeUrl, nil
	}

	if !ok {
		return "", errors.New("Cannot find route by Name " + routeName)
	}
	for k, v := range param {
		routeUrl = strings.Replace(routeUrl, "/:"+k, v, 1)
	}
	return routeUrl, nil
}

type route struct {
	name              string
	bindFunc          func(handlers ...request.HandleFunc)
	handlers          []request.HandleFunc
	relativePath      string
	httpMethod        string
	basicPath         string
	prefixHandlersNum int
	wsHandler         websocket.Handler
}

const maxRouteMapLength = 1000

type engineRoute struct {
	lock sync.RWMutex
	data map[request.EngineHash]chan *route
}

func newEngineRoute() *engineRoute {
	return &engineRoute{
		data: make(map[request.EngineHash]chan *route),
	}
}

func (erm *engineRoute) initEngine(eh request.EngineHash) {
	if erm.data[eh] == nil {
		erm.data[eh] = make(chan *route, maxRouteMapLength)
	}
}

func (erm *engineRoute) Get(eh request.EngineHash) chan *route {
	erm.lock.RLock()
	defer erm.lock.RUnlock()
	return erm.data[eh]
}

func (erm *engineRoute) Set(eh request.EngineHash, r *route) {
	erm.lock.Lock()
	defer erm.lock.Unlock()
	erm.initEngine(eh)
	erm.data[eh] <- r
}

func (erm *engineRoute) Close(eh request.EngineHash) {
	close(erm.data[eh])
}

var engineRouteMap *engineRoute

func init() {
	engineRouteMap = newEngineRoute()
	RouteNameMap = newRouteNameMap()
}

func newRoute(httpMethod string, g *group, relativePath string, bindFunc func(ginHandlers ...request.HandleFunc), handlers ...request.HandleFunc) *route {
	r := route{
		httpMethod:        httpMethod,
		prefixHandlersNum: len(g.RouterGroup.Handlers),
		basicPath:         g.RouterGroup.BasePath(),
		relativePath:      relativePath,
		bindFunc:          bindFunc,
		handlers:          handlers,
	}
	engineRouteMap.Set(g.engineHash, &r)
	return &r
}

func newWsRoute(httpMethod string, g *group, relativePath string, bindWsFunc func(wsHandler websocket.Handler, ginHandlers ...request.HandleFunc), wsHandler websocket.Handler, handlers ...request.HandleFunc) *route {
	r := route{
		wsHandler:         wsHandler,
		httpMethod:        httpMethod,
		prefixHandlersNum: len(g.RouterGroup.Handlers),
		basicPath:         g.RouterGroup.BasePath(),
		relativePath:      relativePath,
		bindFunc: func(handlers ...request.HandleFunc) {
			httpHs := append(handlers, func(context request.Context) {

			})
			bindWsFunc(wsHandler, httpHs...)
		},
		handlers: handlers,
	}
	engineRouteMap.Set(g.engineHash, &r)
	return &r
}

func (r *route) Name(routeName string) policy.RoutePolicier {
	r.name = routeName
	return r
}

func (r *route) Can(policies policy.IPolicy, action policy.Action) {
	r.handlers = append([]request.HandleFunc{
		middleware.Policy(policies, action),
	}, r.handlers...)
}

func (r *route) handlerNum() int {
	return r.prefixHandlersNum + len(r.handlers)
}

func (r *route) absolutePath() string {
	return fmt.Sprintf("%s%s", r.basicPath, r.relativePath)
}

func (r *route) lastHandlerName() string {
	if r.httpMethod == httpMethodWebsocket {
		h := reflect.ValueOf(r.wsHandler).Elem().Type()
		return h.PkgPath() + "." + h.Name()
	}
	if len(r.handlers) > 0 {
		return runtime.FuncForPC(reflect.ValueOf(r.handlers[len(r.handlers)-1]).Pointer()).Name()
	}
	return "nil"
}

func Bind(engine *request.Engine) {
	hash := engine.Hash()
	defer engineRouteMap.Close(hash)

	for r := range engineRouteMap.Get(hash) {
		r.bindFunc(r.handlers...)

		RouteNameMap.set(r.name, r.absolutePath())

		if config.GetMode() != config.ModeProd {
			log.Info(fmt.Sprintf("%-6s %-30s --> %s (%d handlers)\n", r.httpMethod, r.absolutePath(), r.lastHandlerName(), r.handlerNum()))
		}

		if len(engineRouteMap.Get(hash)) <= 0 {
			break
		}
	}

}
