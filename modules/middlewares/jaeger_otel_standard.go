package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetTraceHeaderFromJaeger conver uber-trace-id header to traceparent
// here you can read more about those formats:
//
//	https://www.w3.org/TR/trace-context/#traceparent-header
//	https://www.jaegertracing.io/docs/1.22/client-libraries/#key
func GetTraceHeaderFromJaeger(r *http.Request) {
	uberTraceID := r.Header.Get("Uber-Trace-ID")
	traceparent := r.Header.Get("traceparent")
	if uberTraceID == "" || traceparent != "" {
		return
	}

	traceSlice := strings.Split(uberTraceID, ":")
	if len(traceSlice) < 4 {
		return
	}

	traceID := fmt.Sprintf("%0*s", 32, traceSlice[0])
	spanID := fmt.Sprintf("%0*s", 16, traceSlice[1])
	spanFlags := fmt.Sprintf("%0*s", 2, traceSlice[3])
	traceparent = fmt.Sprintf(
		"00-%s-%s-%s",
		traceID,
		spanID,
		spanFlags,
	)

	r.Header.Set("traceparent", traceparent)
}

// MethodOverride validate authenticity on post
func ConvertJaegerToOtel(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		GetTraceHeaderFromJaeger(c.Request())

		return next(c)
	}
}
