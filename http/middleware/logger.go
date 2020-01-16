package middleware

import (
	"github.com/gofier/framework/request"

	"github.com/gin-gonic/gin"
)

func Logger() request.HandleFunc {
	return func(c request.Context) {
		gin.Logger()(c.GinContext())
	}
}
