package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/junhaideng/go-session"
	"github.com/junhaideng/go-session/client"
)

func BasicAuth() {
	username, passwd := "basic-auth", "password"
	s := session.NewClient(client.WithBasicAuth(username, passwd))
	resp, err := s.Get(fmt.Sprintf("http://httpbin.org/basic-auth/%s/%s", username, passwd))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode == http.StatusOK)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}

type CookiesResponse struct {
	Cookies Cookies `json:"cookies"`
}

type Cookies struct {
	Name string `json:"name"`
}

func Cookie() {
	name, value := "name", "value"
	s := session.NewClient()
	// set cookie
	s.Get(fmt.Sprintf("http://httpbin.org/cookies/set/%s/%s", name, value))

	fmt.Println(s.Cookies())
	// get cookies
	resp, err := s.Get("http://httpbin.org/cookies")

	if err != nil {
		fmt.Println("Get failed: ", err)
		return
	}
	defer resp.Body.Close()
	var res = CookiesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", res)
}

func main() {
	BasicAuth()
	Cookie()
}
