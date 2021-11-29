## Go-Session

### client

Client side session can persist 
- cookies
- basic authorization
- basic proxy authorization

> Note: we persist cookies even when redirection, if you do not want this, specify your own `CheckRedirect` function

Examples can be found in [client.go](examples/client.go)