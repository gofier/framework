package context

import (
	"io"
	"net/http"

	"github.com/gofier/framework/zone"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type IHttpContext interface {
	ILifeCycleContext
	IDataContext

	IRequestParamContext
	IRequestFileContext
	IRequestBindingContext
	IRequestHeaderContext
	IRequestRawContext

	IResponseContext
	IResponseFileContext
	IResponseStreamContext

	Request() *http.Request
	SetRequest(r *http.Request)
	Writer() gin.ResponseWriter
	SetWriter(w gin.ResponseWriter)
	Params() gin.Params
	Accepted() []string
	Keys() map[string]interface{}
	Errors() []*gin.Error

	Deadline() (deadline zone.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type ILifeCycleContext interface {
	Copy() *gin.Context
	HandlerName() string
	HandlerNames() []string
	Handler() gin.HandlerFunc
	Next()
	IsAborted() bool
	Abort()
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, jsonObj interface{})
	AbortWithError(code int, err error) *gin.Error
	Error(err error) *gin.Error
}

type IDataContext interface {
}

type IRequestParamContext interface {
}

type IRequestFileContext interface {
}

type IRequestBindingContext interface {
}

type IRequestHeaderContext interface {
	ClientIP() string
}

type IRequestRawContext interface {
	GetRawData() ([]byte, error)
}

type IResponseContext interface {
	SetAccepted(formats ...string)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Cookie(name string) (string, error)
	Render(code int, r render.Render)
	HTML(code int, name string, obj interface{})
	IndentedJSON(code int, obj interface{})
	SecureJSON(code int, obj interface{})
	JSONP(code int, obj interface{})
	JSON(code int, obj interface{})
	AsciiJSON(code int, obj interface{})
	PureJSON(code int, obj interface{})
	XML(code int, obj interface{})
	YAML(code int, obj interface{})
	ProtoBuf(code int, obj interface{})
	String(code int, format string, values ...interface{})
	Redirect(code int, location string)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
}

type IResponseFileContext interface {
}

type IResponseStreamContext interface {
}
