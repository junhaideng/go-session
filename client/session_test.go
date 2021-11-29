package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	basicAuth := Auth{Username: "basic-auth", Password: "password"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.EscapedPath() == "/basic-auth" {
				u, p, ok := r.BasicAuth()
				if !ok {
					t.Fatal()
				}
				if u != basicAuth.Username && p != basicAuth.Password {
					t.Fatal(u, p)
				}
				t.Log("Check auth succeed")
			}
		}
	}))
	s := New(WithBasicAuth(basicAuth.Username, basicAuth.Password))

	_, err := s.Get(server.URL + "/basic-auth")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", s)
}

func TestProxyAuth(t *testing.T) {
	proxyAuth := Auth{Username: "proxy-auth", Password: "password"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.EscapedPath() == "/proxy-auth" {
				auth := r.Header.Get("Proxy-Authorization")
				if auth == "" {
					t.Fatal("No proxy authorization found")
				}

				a, ok := ParseAuth(auth)
				if !ok {
					t.Fatal()
				}
				if a.Username != proxyAuth.Username && a.Password != proxyAuth.Password {
					t.Fatal("Username or password is not matched")
				}
			}
		}
	}))
	s := New(
		WithProxyAuth(proxyAuth.Username, proxyAuth.Password),
	)

	_, err := s.Get(server.URL + "/proxy-auth")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v\n", s)
}

func TestCookie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			if r.URL.EscapedPath() == "/set-cookie" {
				w.Header().Set("Set-Cookie", "name=test")
				w.WriteHeader(http.StatusOK)
			}
			if r.URL.EscapedPath() == "/check-cookie" {
				c, err := r.Cookie("name")
				if err != nil {
					t.Fatal(err)
				}
				if c.Value != "test" {
					t.Fatal("Cookie is not matched")
				}
				t.Log("Check cookie succeed")
			}
		}
	}))
	s := New()

	// -------------------check cookie-------------------------
	_, err := s.Get(server.URL + "/set-cookie")
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Get(server.URL + "/check-cookie")
	if err != nil {
		t.Fatal(err)
	}

	// -----------------------------------------------------------
	t.Logf("%#v\n", s)
}

func TestRepeatCookie(t *testing.T) {
	var count = 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.EscapedPath() == "/set-cookie" {
				count++
				w.Header().Set("Set-Cookie", fmt.Sprintf("name=%d", count))
				w.WriteHeader(http.StatusOK)
			}
			if r.URL.EscapedPath() == "/check-cookie" {
				c, err := r.Cookie("name")
				if err != nil {
					t.Fatal(err)
				}
				if c.Value != fmt.Sprintf("%d", count) {
					t.Fatal("Cookie is not matched")
				}
				t.Log("Check cookie succeed")
			}
		}
	}))
	s := New()

	// -------------------check cookie-------------------------
	_, err := s.Get(server.URL + "/set-cookie")
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Get(server.URL + "/set-cookie")
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Get(server.URL + "/check-cookie")
	if err != nil {
		t.Fatal(err)
	}

	// -----------------------------------------------------------
	t.Logf("%#v\n", s)
}

func TestCheckRedirectCookie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.EscapedPath() == "/set-cookie" {
				w.Header().Set("Set-Cookie", "name=value")
				http.Redirect(w, r, "/redirect-1", http.StatusFound)
			}

			if r.URL.EscapedPath() == "/redirect-1" {
				http.Redirect(w, r, "/redirect-2", http.StatusFound)
			}
			if r.URL.EscapedPath() == "/redirect-2" {
				w.WriteHeader(http.StatusOK)
			}

			if r.URL.EscapedPath() == "/check-cookie" {
				c, err := r.Cookie("name")
				if err != nil {
					t.Fatal(err)
				}
				if c.Value != "value" {
					t.Fatal("Cookie is not matched")
				}
				t.Log("Check cookie succeed")
			}
		}
	}))
	s := New()

	// -------------------check cookie-------------------------
	_, err := s.Get(server.URL + "/set-cookie")
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Get(server.URL + "/check-cookie")
	if err != nil {
		t.Fatal(err)
	}

	// -----------------------------------------------------------
	t.Logf("%#v\n", s)
}
