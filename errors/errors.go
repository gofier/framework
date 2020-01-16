package errors

import (
	"fmt"

	"github.com/gofier/framework/log"

	"github.com/ztrue/tracerr"
)

func ErrPrintln(err error, fields map[string]interface{}) {
	startFrom := 2
	if err != nil {
		return
	}
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	if startFrom > len(frameList) || len(frameList)-2 <= 0 {
		log.ErrorWithFields(err.Error(), fields)
	}

	traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	traceErr = tracerr.CustomError(err, frameList)

	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["gofer_trace"] = tracerr.SprintSource(traceErr, 0)
	log.ErrorWithFields(err.Error(), fields)
}

func ErrPrint(err error, startFrom int, fields map[string]interface{}) string {
	if err != nil {
		return ""
	}
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	traceErr = tracerr.CustomError(err, frameList)

	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["gofer_trace"] = tracerr.SprintSource(traceErr)
	return fmt.Sprint(err.Error(), fields)
}
