package policy

import (
	"github.com/gofier/framework/context"
	"github.com/gofier/framework/request/http/auth"
)

type IPolicyContext interface {
	context.ILifeCycleContext
	context.IResponseContext
	auth.IAuthContext
	auth.RequestIUser
}
