package middleware

import (
	"github.com/gofier/framework/policy"
	"github.com/gofier/framework/request"
)

func Policy(_policy policy.IPolicy, action policy.Action) request.HandleFunc {
	return func(c request.Context) {
		policy.Middleware(_policy, action, c, c.Params())
	}
}
