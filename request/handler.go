package request

import (
	"github.com/gofier/framework/request/http"

	"github.com/gin-gonic/gin"
)

func ConvertHandlers(handlers []HandleFunc) (ginHandlers []gin.HandlerFunc) {
	for _, h := range handlers {
		handler := h
		ginHandlers = append(ginHandlers, func(c *gin.Context) {
			goferCtx := http.ConvertContext(c)
			handler(goferCtx)
		})
	}
	return
}
