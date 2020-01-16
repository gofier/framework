package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gofier/framework/config"
	"github.com/gofier/framework/log"
	"github.com/gofier/framework/request"
	"github.com/gofier/framework/zone"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() request.HandleFunc {
	return func(c request.Context) {
		if config.GetBoolean("app.debug", true) {
			startedAt := zone.Now()
			requestHeader := c.Request().Header
			requestData, err := c.GetRawData()
			if err != nil {
				log.Error(err)
				c.Next()
			}
			r := c.Request()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(requestData))
			c.SetRequest(r)

			responseWriter := &responseWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer(),
			}
			c.SetWriter(responseWriter)

			defer log.InfoWithFields(c.ClientIP(), map[string]interface{}{
				"Method":         c.Request().Method,
				"Path":           c.Request().RequestURI,
				"Proto":          c.Request().Proto,
				"Status":         responseWriter.Status(),
				"UA":             c.Request().UserAgent(),
				"Latency":        zone.Now().Sub(startedAt),
				"RequestHeader":  requestHeader,
				"RequestBody":    string(requestData),
				"ResponseHeader": responseWriter.Header(),
				"ResponseBody":   responseWriter.body.String(),
			})
		}

		c.Next()
	}
}
