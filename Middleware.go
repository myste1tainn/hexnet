package msnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/myste1tainn/hexlog"
	"github.com/myste1tainn/msfnd"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type TraceOptions struct {
	SkipPaths []string
}

func TraceWithDefaultOptions() gin.HandlerFunc {
	return Trace(TraceOptions{
		SkipPaths: []string{"/builds", "/health", "/metrics"},
	})
}

func newLogger(ctx *gin.Context) *log.Logger {
	l := log.L.NewChildLogger()

	traceId, ok := l.DefaultPayload[GcpJsonLoggingTraceKey].(string)
	if !ok || traceId == "" {
		l.DefaultPayload[GcpJsonLoggingTraceKey] = ctx.Request.Header.Get(HttpRequestHeaderTraceKey)
		traceId, ok := l.DefaultPayload[GcpJsonLoggingTraceKey].(string)
		if !ok || traceId == "" {
			id := strings.ReplaceAll(uuid.New().String(), "-", "")
			l.DefaultPayload[GcpJsonLoggingTraceKey] = id
			l.DefaultPayload[HttpRequestHeaderTraceKey] = id
		}
	}

	spanId, ok := l.DefaultPayload[GcpJsonLoggingSpanIdKey].(string)
	if !ok || spanId == "" {
		l.DefaultPayload[GcpJsonLoggingSpanIdKey] = ctx.Request.Header.Get(HttpRequestHeaderSpanIdKey)
	}

	return l
}

func setLogTraceToRouteContext(ctx *gin.Context, l *log.Logger) {
	v, ok := ctx.Get(msfnd.KeyRouteContext)
	if !ok || v == nil {
		return
	}

	rctx, ok := v.(*msfnd.RouteContext)
	if !ok || rctx == nil {
		rctx = newRouteContext(ctx)
	}

	if traceId, ok := l.DefaultPayload[GcpJsonLoggingTraceKey].(string); ok && traceId != "" {
		rctx.SetTrace(traceId)
	}
	if spanId, ok := l.DefaultPayload[GcpJsonLoggingSpanIdKey].(string); ok && spanId != "" {
		rctx.SetSpanId(spanId)
	}
}

func newRouteContext(ctx *gin.Context) *msfnd.RouteContext {
	rctx := &msfnd.RouteContext{}
	ctx.Set(msfnd.KeyRouteContext, rctx)
	return rctx
}

func Trace(opts TraceOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(opts.SkipPaths, ctx.Request.URL.Path) {
			return
		}

		l := newLogger(ctx)
		setLoggerToContext(l, ctx)
		defer l.Destroy()

		setLogTraceToRouteContext(ctx, l)

		var reqBody []byte
		if ctx.Request.Body != nil {
			reqBody, _ = ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		headersMap := map[string]any{}
		for k, v := range ctx.Request.Header {
			if len(v) == 1 {
				headersMap[k] = fmt.Sprintf("%s", v[0])
			} else {
				headersMap[k] = fmt.Sprintf("%v", v)
			}
		}

		// Get the body
		var m map[string]any
		var a []any
		var b any
		err := json.Unmarshal(reqBody, &m)
		if err != nil {
			err := json.Unmarshal(reqBody, &a)
			if err != nil {
				b = a
			} else {
				b = ""
			}
		} else {
			b = m
		}

		p := RequestLogPayload{
			Type: "http-inbound-request",
			HttpRequest: RequestPayload{
				Headers:   headersMap,
				Host:      ctx.Request.Host,
				Path:      ctx.Request.RequestURI,
				Method:    ctx.Request.Method,
				Body:      b,
				Timestamp: time.Now(),
			},
		}
		l.InfoJsonf(p, "BEGIN | INBOUND | %7s | %s%s", ctx.Request.Method, ctx.Request.Host, ctx.Request.RequestURI)

		//ctx.SetEventHandler(app.EventHandler{
		//	OnBeforeResponse: func(code int, headers map[string]any, res any) {
		//		defer l.Destroy()
		//		if l == nil {
		//			l = log.L
		//			l.Warnf("logger was deinit before it is used for END INBOUND response logging")
		//		}
		//		payload := ResponseLogPayload{
		//			Type: "http-inbound-response",
		//			HttpResponse: ResponsePayload{
		//				Code:      ctx.Status(),
		//				Headers:   headers,
		//				Host:      ctx.Request.Host,
		//				Path:      ctx.Request.RequestURI,
		//				Method:    ctx.Request.Method,
		//				Body:      getBody(res),
		//				Timestamp: time.Now(),
		//			},
		//		}
		//		if ctx.Status() > 399 || ctx.Status() < 200 {
		//			l.ErrorJsonf(payload, "END   | INBOUND | %-7d | %7s | %s | %s%s", ctx.Status(), ctx.Request().Method, time.Since(p.HttpRequest.Timestamp).String(), ctx.Request().Host, ctx.Request().RequestURI)
		//		} else {
		//			l.InfoJsonf(payload, "END   | INBOUND | %-7d | %7s | %s | %s%s", ctx.Status(), ctx.Request().Method, time.Since(p.HttpRequest.Timestamp).String(), ctx.Request().Host, ctx.Request().RequestURI)
		//		}
		//	},
		//})

		ctx.Next()

	}
}
