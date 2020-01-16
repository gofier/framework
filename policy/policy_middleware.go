package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserNotPermitError struct {
}

func (e UserNotPermitError) Error() string {
	return "user has no permission"
}

func Middleware(policy IPolicy, action Action, c IPolicyContext, params []gin.Param) {
	routeParamMap := make(map[string]string)
	for _, param := range params {
		routeParamMap[param.Key] = param.Value
	}

	if err := c.ScanUser(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"error": err})
		c.Abort()
		return
	}

	c.Next()
}

func forbid(c IPolicyContext) {
	c.JSON(http.StatusForbidden, map[string]interface{}{
		"error": UserNotPermitError{}.Error(),
	})
}
