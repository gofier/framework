package request

import (
	"github.com/gofier/framework/context"
	"github.com/gofier/framework/request/http/auth"
	"github.com/gofier/framework/utils/jwt"

	"github.com/gin-gonic/gin"
)

type Context interface {
	context.IHttpContext

	GinContext() *gin.Context

	SetAuthClaim(claims *jwt.UserClaims)
	SetIUserModel(iUser auth.IUser)

	auth.IAuthContext
	auth.RequestIUser
}
