package client

import (
	"encoding/base64"
	"strings"
	"unsafe"
)

// Auth is a simple basic auth using base64 encoding
type Auth struct {
	Username string
	Password string
	result   string
}

func (a *Auth) encode() string {
	if a.result != "" {
		return a.result
	}
	tmp := a.Username + ":" + a.Password
	a.result = "Basic " + base64.StdEncoding.EncodeToString(*(*[]byte)(unsafe.Pointer(&tmp)))
	return a.result
}

func ParseAuth(auth string) (a *Auth, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return &Auth{
		Username: cs[:s],
		Password: cs[s+1:],
	}, true
}
