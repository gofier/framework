package middleware

import (
	"github.com/gofier/framework/request"

	"github.com/gin-gonic/gin"
)

func Recovery() request.HandleFunc {
	return func(c request.Context) {
		gin.Recovery()(c.GinContext())
	}
}
