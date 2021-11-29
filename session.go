package session

import "github.com/junhaideng/go-session/client"

func NewClient(options ...client.Option) client.SessionClient{
	return client.New(options...)
}