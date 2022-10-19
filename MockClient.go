package msnet

import (
	"crypto/tls"

	"github.com/imroc/req/v3"
	log "github.com/myste1tainn/hexlog"
	"github.com/myste1tainn/msfnd"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	*mock.Mock
}

// DisableInsecureSkipVerify implements Client
func (c *MockClient) DisableInsecureSkipVerify() Client {
	return c
}

// EnableInsecureSkipVerify implements Client
func (c *MockClient) EnableInsecureSkipVerify() Client {
	return c
}

// NewRequest implements Client
func (c *MockClient) NewRequest(rctx *msfnd.RouteContext, fp string, api Api) Request {
	return &DefaultRequest{
		Request:      &req.Request{},
		logger:       log.L,
		routeContext: rctx,
		fullPath:     fp,
		api:          api,
	}
}

// Request implements Client
func (c *MockClient) Request(config ConfigProtocol, apiId string) Request {
	api := config.GetApis()[apiId]
	return c.NewRequest(&msfnd.RouteContext{}, config.FullPath(api), api)
}

// RequestWithContext implements Client
func (c *MockClient) RequestWithContext(rctx *msfnd.RouteContext, config ConfigProtocol, apiId string) Request {
	api := config.GetApis()[apiId]
	return c.NewRequest(rctx, config.FullPath(api), api)
}

// SetCertFromFile implements Client
func (c *MockClient) SetCertFromFile(certFile string, keyFile string) Client {
	return c
}

// SetCerts implements Client
func (c *MockClient) SetCerts(certs ...tls.Certificate) Client {
	return c
}

// WithLogger implements Client
func (c *MockClient) WithLogger(logger *log.Logger) Client {
	return c
}

func NewMockClient() *MockClient {
	return &MockClient{
		Mock: &mock.Mock{},
	}
}
