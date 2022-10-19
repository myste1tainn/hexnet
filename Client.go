package msnet

import (
	"crypto/tls"

	"github.com/imroc/req/v3"
	log "github.com/myste1tainn/hexlog"
	"github.com/myste1tainn/msfnd"
)

type Client interface {
	SetCertFromFile(certFile, keyFile string) Client
	SetCerts(certs ...tls.Certificate) Client
	WithLogger(logger *log.Logger) Client
	EnableInsecureSkipVerify() Client
	DisableInsecureSkipVerify() Client
	RequestWithContext(rctx *msfnd.RouteContext, config ConfigProtocol, apiId string) Request
	Request(config ConfigProtocol, apiId string) Request
	NewRequest(rctx *msfnd.RouteContext, fp string, api Api) Request
}

type DefaultClient struct {
	*req.Client
	config ClientConfig
	logger *log.Logger
}

var isGlobalDevMode = false

func DevMode() {
	req.DevMode()
	isGlobalDevMode = true
}

func NewClient(config ClientConfig) Client {
	return &DefaultClient{
		Client: req.C(),
		config: config,
		logger: log.L,
	}
}

func (c *DefaultClient) SetCertFromFile(certFile string, keyFile string) Client {
	c.Client.SetCertFromFile(certFile, keyFile)
	return c
}

func (c *DefaultClient) SetCerts(certs ...tls.Certificate) Client {
	c.Client.SetCerts(certs...)
	return c
}

func (c *DefaultClient) WithLogger(logger *log.Logger) Client {
	c.logger = logger
	return c
}

func (c *DefaultClient) getLogger() *log.Logger {
	if c.logger == nil {
		return log.L
	} else {
		return c.logger
	}
}

func (c *DefaultClient) RequestWithContext(rctx *msfnd.RouteContext, config ConfigProtocol, apiId string) Request {
	l := c.getLogger()

	api := config.GetApis()[apiId]
	if api == (Api{}) {
		l.Panicf("trying to get unknown API spec = %s", apiId)
	}

	fp := config.FullPath(api)
	if fp == "" {
		l.Panicf("got empty url while trying to construct full path for API spec = %v", api)
	}

	l.Debugf("making request with apiId = %s", apiId)
	return c.NewRequest(rctx, fp, api)
}

func (c *DefaultClient) Request(config ConfigProtocol, apiId string) Request {
	return c.RequestWithContext(nil, config, apiId)
}

func (c *DefaultClient) NewRequest(rctx *msfnd.RouteContext, fp string, api Api) Request {
	return &DefaultRequest{
		Request:      c.R(),
		logger:       c.getLogger(),
		routeContext: rctx,
		fullPath:     fp,
		api:          api,
	}
}

func (c *DefaultClient) EnableInsecureSkipVerify() Client {
	c.Client.EnableInsecureSkipVerify()
	return c
}

func (c *DefaultClient) DisableInsecureSkipVerify() Client {
	c.Client.DisableInsecureSkipVerify()
	return c
}
