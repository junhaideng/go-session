// Package client implements client side session
// which persists cookies, basic authorization and basic proxy authorization
package client

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SessionClient works just like a http client
type SessionClient interface {
	SetBasicAuth(username, password string)
	SetProxyAuth(username, password string)
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	PostForm(url string, data url.Values) (*http.Response, error)
	Head(url string) (*http.Response, error)
	Cookies() []*http.Cookie
}

// keep cookie when redirect
var defaultCheckRedirect = func(c *client) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		for _, r := range via {
			if r.Response != nil {
				for _, cookie := range r.Response.Cookies() {
					c.cookies[cookie.Name] = cookie
				}
			}
		}

		if req.Response != nil {
			for _, cookie := range req.Response.Cookies() {
				c.cookies[cookie.Name] = cookie
			}
		}
		return nil
	}
}

func New(options ...Option) SessionClient {
	c := &client{cookies: make(map[string]*http.Cookie)}
	for _, option := range options {
		option(c)
	}
	if c.checkRedirect == nil {
		c.checkRedirect = defaultCheckRedirect(c)
	}
	if c.client == nil {
		c.client = &http.Client{
			CheckRedirect: c.checkRedirect,
		}
	}
	return c
}

type client struct {
	client        *http.Client
	proxyAuth     *Auth
	basicAuth     *Auth
	cookies       map[string]*http.Cookie
	checkRedirect func(req *http.Request, via []*http.Request) error
}

// SetBasicAuth set authorization
func (c *client) SetBasicAuth(username, password string) {
	c.basicAuth = &Auth{
		Username: username,
		Password: password,
	}
}

func (c *client) SetProxyAuth(username, password string) {
	c.proxyAuth = &Auth{
		Username: username,
		Password: password,
	}
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func (c *client) Get(url string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.do(req)
}

func (c *client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *client) Head(url string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	return c.do(req)
}

func (c *client) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// add cookie
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// add authorization header
	if c.basicAuth != nil {
		req.Header.Set("Authorization", c.basicAuth.encode())
	}

	if c.proxyAuth != nil {
		req.Header.Set("Proxy-Authorization", c.proxyAuth.encode())
	}
	return req, err
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// keep the most fresh cookie
	for _, cookie := range resp.Cookies() {
		c.cookies[cookie.Name] = cookie
	}
	return resp, err
}

func (c *client) Cookies() []*http.Cookie {
	res := make([]*http.Cookie, 0, len(c.cookies))
	for _, cookie := range c.cookies {
		res = append(res, cookie)
	}
	return res
}
