## Go-Session
[![Go](https://github.com/junhaideng/go-session/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/junhaideng/go-session/actions/workflows/go.yml)

### client

Client side session can persist 
- cookies
- basic authorization
- basic proxy authorization

> Note: we persist cookies even when redirection, if you do not want this, specify your own `CheckRedirect` function

Examples can be found in [client.go](examples/client.go)