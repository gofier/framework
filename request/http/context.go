package http

import (
	"net/http"
	"reflect"
	"unicode/utf8"

	"github.com/gofier/framework/config"
	"github.com/gofier/framework/request/http/auth"
	"github.com/gofier/framework/utils/jwt"

	"github.com/gin-gonic/gin"
)

const CONTEXT_CLAIM_KEY = "GOFIER_CONTEXT_CLAIM"
const CONTEXT_IUSER_MODEL_KEY = "GOFIER_CONTEXT_IUSER_MODEL"

type httpContext struct {
	*gin.Context
	*auth.RequestUser
}

func (c *httpContext) GinContext() *gin.Context {
	return c.Context
}

func (c *httpContext) Request() *http.Request {
	return c.Context.Request
}

func (c *httpContext) Writer() gin.ResponseWriter {
	return c.Context.Writer
}

func (c *httpContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}

func (c *httpContext) SetWriter(w gin.ResponseWriter) {
	c.Context.Writer = w
}

func (c *httpContext) Params() gin.Params {
	return c.Context.Params
}

func (c *httpContext) Accepted() []string {
	return c.Context.Accepted
}

func (c *httpContext) Keys() map[string]interface{} {
	return c.Context.Keys
}

func (c *httpContext) Errors() []*gin.Error {
	return c.Context.Errors
}

func (c *httpContext) SetAuthClaim(claims *jwt.UserClaims) {
	c.Set(CONTEXT_CLAIM_KEY, claims)
}

func (c *httpContext) AuthClaimID() (ID uint, exist bool) {
	claims, exist := c.Get(CONTEXT_CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}

func (c *httpContext) SetIUserModel(iuser auth.IUser) {
	c.Set(CONTEXT_IUSER_MODEL_KEY, iuser)
}

func (c *httpContext) IUserModel() auth.IUser {
	iuser, exist := c.Get(CONTEXT_IUSER_MODEL_KEY)

	var typeof reflect.Type
	if !exist {
		typeof = reflect.TypeOf(config.GetInterface("auth.model_ptr"))
	} else {
		typeof = reflect.TypeOf(iuser.(auth.IUser))
	}

	ptr := reflect.New(typeof).Elem()
	val := reflect.New(typeof).Elem()
	ptr.Set(val)
	return ptr.Interface().(auth.IUser)
}

func (c *httpContext) ScanUserWithJSON() (isAbort bool) {
	if err := c.ScanUser(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": err,
		})
		return true
	}
	return false
}

func ConvertContext(c *gin.Context) *httpContext {
	_c := &httpContext{
		Context:     c,
		RequestUser: &auth.RequestUser{},
	}
	_c.RequestUser.SetContext(_c)
	return _c
}
