package blb

import (
	"github.com/drinktee/bce-sdk-go/bce"
)

// Endpoint contains all endpoints of Baidu Cloud BCC.
var Endpoint = map[string]string{
	"bj": "blb.bj.baidubce.com",
	"gz": "blb.gz.baidubce.com",
}

// Config contains all options for bos.Client.
type Config struct {
	*bce.Config
}

func NewConfig(config *bce.Config) *Config {
	return &Config{config}
}

// Client is the bos client implemention for Baidu Cloud BOS API.
type Client struct {
	*bce.Client
}

func NewBLBClient(config *Config) *Client {
	bceClient := bce.NewClient(config.Config)
	return &Client{bceClient}
}

// GetURL generates the full URL of http request for Baidu Cloud BOS API.
func (c *Client) GetURL(version string, params map[string]string) string {
	host := c.Endpoint
	if host == "" {
		host = Endpoint[c.GetRegion()]
	}
	uriPath := version
	return c.Client.GetURL(host, uriPath, params)
}
