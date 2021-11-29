package client

import (
	"net/http"
)

// Option using for set parameters
type Option func(*client)

func WithClient(cli *http.Client) Option {
	return func(c *client) {
		c.client = cli
	}
}

func WithBasicAuth(username, password string) Option {
	return func(c *client) {
		c.SetBasicAuth(username, password)
	}
}

func WithProxyAuth(username, password string) Option {
	return func(c *client) {
		c.SetProxyAuth(username, password)
	}
}

func WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) Option{
	return func(c *client) {
		c.checkRedirect = checkRedirect
	}
}