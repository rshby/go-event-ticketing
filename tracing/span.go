package tracing

import (
	"context"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Start starts a span
func Start(ctx context.Context) (context.Context, trace.Span) {
	var spanName = "unknow_function"
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			spanName = fn.Name()
		}
	}

	if gCtx, ok := ctx.(*gin.Context); ok {
		newCtx, span := otel.Tracer(config.OtlpServiceName()).Start(gCtx.Request.Context(), spanName)
		gCtx.Request = gCtx.Request.WithContext(newCtx)
		return newCtx, span
	}

	return otel.Tracer(config.OtlpServiceName()).Start(ctx, spanName)
}
