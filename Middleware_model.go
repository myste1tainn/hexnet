package msnet

import (
	"time"
)

type RequestLogPayload struct {
	Type        string         `json:"type"`
	HttpRequest RequestPayload `json:"httpRequest"`
}

type RequestPayload struct {
	Headers   map[string]any `json:"headers"`
	Host      string         `json:"host"`
	Path      string         `json:"path"`
	Method    string         `json:"method"`
	Body      any            `json:"body"`
	Timestamp time.Time      `json:"timestamp"`
}

type ResponseLogPayload struct {
	Type         string          `json:"type"`
	HttpResponse ResponsePayload `json:"httpResponse"`
}

type ResponsePayload struct {
	Code      int            `json:"code"`
	Headers   map[string]any `json:"headers"`
	Host      string         `json:"host"`
	Path      string         `json:"path"`
	Method    string         `json:"method"`
	Body      any            `json:"body"`
	Timestamp time.Time      `json:"timestamp"`
}
