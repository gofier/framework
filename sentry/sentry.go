package sentry

import (
	"fmt"

	"github.com/gofier/framework/config"

	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

var enabled bool

func Init() {
	enabled = config.GetBoolean("sentry.enable")
	if enabled {
		err := raven.SetDSN(fmt.Sprintf("https://%s:%s@%s/%s",
			config.GetString("sentry.key"),
			config.GetString("sentry.secret"),
			config.GetString("sentry.host"),
			config.GetString("sentry.project"),
		))
		if err != nil {
			panic(err)
		}
	}
}

func Use(r *gin.Engine, onlySendOnCrash bool) {
	if enabled {
		r.Use(sentry.Recovery(raven.DefaultClient, onlySendOnCrash))
	}
}

func CaptureError(err error) {
	if enabled {
		raven.CaptureErrorAndWait(err, map[string]string{
			"env": config.GetString("app.env"),
		})
	}
}

func CaptureMsg(msg string, field map[string]interface{}) {
	if enabled {
		raven.CaptureMessage(fmt.Sprintf("%s - %v", msg, field), map[string]string{
			"env": config.GetString("app.env"),
		})
	}
}

func CapturePanic(handler func()) {
	if enabled {
		raven.CapturePanic(handler, map[string]string{
			"env": config.GetString("app.env"),
		})
	}
}
