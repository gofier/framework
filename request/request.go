package request

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type HandleFunc func(Context)

type Engine struct {
	*gin.Engine
}

func New() *Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		//
	}
	return &Engine{
		Engine: gin.New(),
	}
}

type EngineHash = string

func (e *Engine) Hash() EngineHash {
	return fmt.Sprintf("%x", e)
}

func (e *Engine) Use(handlers ...HandleFunc) {
	e.Engine.Use(ConvertHandlers(handlers)...)
}

func (e *Engine) UseGin(handlerFunc ...gin.HandlerFunc) {
	e.Engine.Use(handlerFunc...)
}

func (e *Engine) GinEngine() *gin.Engine {
	return e.Engine
}
