package msnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	log "github.com/myste1tainn/hexlog"
	"github.com/myste1tainn/msfnd"
)

type Request interface {
	Call(result any, error any) (*req.Response, error)
	SetHeader(key, value string) Request
	SetQueryParam(key, value string) Request
	SetPathParams(params map[string]string) Request
	SetPathParam(key string, value string) Request
	SetBody(body any) Request
	SetBearerAuthToken(value string) Request
	SetContentType(value string) Request
	EnableDump() Request
	EnableTrace() Request
	DisableTrace() Request
	AddQueryParam(key, value string) Request
}
type DefaultRequest struct {
	*req.Request
	logger       *log.Logger
	routeContext *msfnd.RouteContext
	fullPath     string
	api          Api
	bodyObject   any
}

func (r *DefaultRequest) getLogger() *log.Logger {
	if r.logger == nil {
		return log.L
	} else {
		return r.logger
	}
}

func (r *DefaultRequest) autoPassHeaders() {
	l := r.getLogger()
	if r.routeContext != nil {
		defer func() {
			if r.routeContext.Authorization != "" {
				r.SetHeader("X-Authorization", r.routeContext.Authorization)
			}
		}()
		var m map[string]any
		dat, err := json.Marshal(r.routeContext)
		if err != nil {
			l.Warnf("[msnet] auto pass headers cannot be performed, err = %v", err)
			return
		} else if err := json.Unmarshal(dat, &m); err != nil {
			l.Warnf("[msnet] auto pass headers cannot be performed, err = %v", err)
			return
		} else {
			for k, v := range m {
				switch t := v.(type) {
				case string:
					r.SetHeader(k, t)
				default:
					r.SetHeader(k, fmt.Sprintf("%v", t))
				}
			}
		}
	} else {
		l.Warnf("[msnet] auto pass headers cannot be performed, route context = nil")
	}
}

func (r *DefaultRequest) Call(resultRes any, errRes any) (*req.Response, error) {
	if resultRes != nil {
		r.SetResult(resultRes)
	}
	if errRes != nil {
		r.SetError(errRes)
	}
	r.autoPassHeaders()

	r.putBeginLog()
	startTime := time.Now()
	var res *req.Response
	var err error
	switch strings.ToLower(r.api.Method) {
	case "get":
		res, err = r.Get(r.fullPath)
	case "post":
		res, err = r.Post(r.fullPath)
	case "put":
		res, err = r.Put(r.fullPath)
	case "patch":
		res, err = r.Patch(r.fullPath)
	case "delete":
		res, err = r.Delete(r.fullPath)
	case "head":
		res, err = r.Head(r.fullPath)
	case "options":
		res, err = r.Options(r.fullPath)
	default:
		res, err = nil, errors.New("unknown http method = "+r.api.Method)
	}
	r.putLog(res, err, resultRes, errRes, startTime)
	return res, err
}

func httpHeaderToMap(h http.Header) map[string]any {
	m := map[string]any{}
	for k, v := range h {
		m[k] = strings.Join(v, ",")
	}
	return m
}

func (r *DefaultRequest) putBeginLog() {
	l := r.getLogger()

	requestLogPayload := RequestLogPayload{
		Type: "http-outbound-request",
		HttpRequest: RequestPayload{
			Headers:   httpHeaderToMap(r.Headers),
			Host:      r.fullPath,
			Path:      r.fullPath,
			Method:    strings.ToUpper(r.api.Method),
			Body:      r.bodyObject,
			Timestamp: time.Now(),
		},
	}

	l.InfoJsonf(requestLogPayload, "BEGIN | OUTBOUND | %7s | %s", strings.ToUpper(r.api.Method), r.fullPath)
}

func (r *DefaultRequest) putLog(res *req.Response, err error, resultRes any, errRes any, startTime time.Time) {
	l := r.getLogger()
	h := map[string]any{}
	if res != nil && res.Response != nil {
		h = httpHeaderToMap(res.Response.Header)
	}

	var createLogPayload = func(body any) ResponseLogPayload {
		return ResponseLogPayload{
			Type: "http-outbound-response",
			HttpResponse: ResponsePayload{
				Headers:   h,
				Host:      r.fullPath,
				Path:      r.fullPath,
				Method:    strings.ToUpper(r.api.Method),
				Body:      body,
				Timestamp: time.Now(),
			},
		}
	}

	if err != nil {
		l.ErrorJsonf(createLogPayload(nil), "END   | OUTBOUND | %-7d | %7s | %s | %s | %s", getStatusCode(res), strings.ToUpper(r.api.Method), time.Since(startTime).String(), r.fullPath, err)
	} else if res.IsError() {
		l.ErrorJsonf(createLogPayload(errRes), "END   | OUTBOUND | %-7d | %7s | %s |%s | %s", getStatusCode(res), strings.ToUpper(r.api.Method), time.Since(startTime).String(), r.fullPath, err)
	} else {
		l.InfoJsonf(createLogPayload(resultRes), "END   | OUTBOUND | %-7d | %7s | %s |%s | %s", getStatusCode(res), strings.ToUpper(r.api.Method), time.Since(startTime).String(), r.fullPath, err)
	}
}

func getStatusCode(res *req.Response) int {
	if res == nil {
		return -1
	}

	if res.Response == nil {
		return -1
	}

	return res.StatusCode
}

func (r *DefaultRequest) SetHeader(key, value string) Request {
	r.Request.SetHeader(key, value)
	return r
}
func (r *DefaultRequest) AddQueryParam(key, value string) Request {
	r.Request.AddQueryParam(key, value)
	return r
}
func (r *DefaultRequest) SetQueryParam(key, value string) Request {
	r.Request.SetQueryParam(key, value)
	return r
}
func (r *DefaultRequest) SetPathParams(params map[string]string) Request {
	r.Request.SetPathParams(params)
	return r
}
func (r *DefaultRequest) SetPathParam(key string, value string) Request {
	r.Request.SetPathParam(key, value)
	return r
}
func (r *DefaultRequest) SetBody(body any) Request {
	r.bodyObject = body
	r.Request.SetBody(body)
	return r
}
func (r *DefaultRequest) SetBearerAuthToken(value string) Request {
	r.Request.SetBearerAuthToken(value)
	return r
}
func (r *DefaultRequest) SetContentType(value string) Request {
	r.Request.SetContentType(value)
	return r
}
func (r *DefaultRequest) EnableDump() Request {
	r.Request.EnableDump()
	return r
}
func (r *DefaultRequest) EnableTrace() Request {
	r.Request.EnableTrace()
	return r
}
func (r *DefaultRequest) DisableTrace() Request {
	r.Request.DisableTrace()
	return r
}
