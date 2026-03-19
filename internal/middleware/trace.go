package middleware

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx      = otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
			spanName = fmt.Sprintf("[%s] %s", strings.ToUpper(c.Request.Method), c.Request.URL.String())
		)

		ctx, span := tracing.StartWithName(c.Request.Context(), spanName)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		startTime := time.Now()

		// execute controller
		c.Next()

		duration := time.Since(startTime)

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.client_ip", c.ClientIP()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.String("http.request.body", string(reqBodyBytes)),
			attribute.String("http.response.body", blw.body.String()),
			attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
		)

		if len(c.Errors) > 0 {
			span.RecordError(fmt.Errorf(c.Errors.String()))
			span.SetStatus(codes.Error, "Request failed")
		} else if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP Error %d", c.Writer.Status()))
		}
	}
}
