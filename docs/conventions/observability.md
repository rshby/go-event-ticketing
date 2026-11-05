# Observability Conventions

## Tracing

- **Setup:** `tracing.ConnectOTLPTrace` (`tracing/common.go`) runs once in the `server` command. The service name is `OTLP_SERVICE_NAME`.
- **Root span:** `middleware.TraceMiddleware` (`internal/middleware/trace.go`), named `"[METHOD] /url"`. It adds `http.method`, `http.url`, `http.client_ip`, `http.user_agent`, `http.status_code`, `http.request.body`, `http.response.body` and `http.duration_ms`. The span status is Error when `c.Errors` is not empty or status ≥ 400.
- **Child spans:** every controller, service, repository and cacher method calls `tracing.Start(ctx)`. The span name is the caller's full function name, found with `runtime.Caller`. Use `tracing.StartWithName(ctx, name)` only when the function name isn't a useful name.
- **`*gin.Context`:** `tracing.Start(c)` reads and replaces `c.Request`'s context. In controllers pass `c` in and use the returned `ctx` from then on.
- **DB:** the GORM otel plugin adds a span for each query (metrics off).
- **Errors:** to mark a span failed at the point where you handle an error, use `span.RecordError(err)` and `span.SetStatus(codes.Error, …)`.

*As built:* incoming `traceparent` is ignored ([KI-09](../known-issues.md#ki-09)).

## Personal data

Request and response bodies are copied onto the root span. Before any order endpoint ships, the middleware must skip or redact bodies for order routes ([KI-22](../known-issues.md#ki-22), [RULE-TKT-09](../invariants/ticketing.md)).

## Logging

- logrus to stdout, default text formatter.
- Per-method logger: `logrus.WithContext(ctx).WithFields(logrus.Fields{"context": helper.DumpIncomingContext(ctx), ...})`. Add the method's key inputs as fields (`"id"`, `"ticket_id"`). `helper.Dump(&dto)` is OK for non-personal DTOs only.
- Levels: `Error` for failures, `Info` for meaningful success ("success reserve slot"), no `Debug` noise.
- Log an error once, where it is first seen or handled ([golang.md](golang.md#method-boilerplate)).
  *As built:* every layer logs the same error again on the way up. Don't extend this in new code.
