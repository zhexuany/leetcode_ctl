package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	LEETCODE_BASE_URL = "https://leetcode.com"
)

type Client interface {
	Ping(timeout time.Duration) (time.Duration, string, error)
	Login() error
	Logout() error
	Search(interface{}) error
	Run(path string) error
	Submit(path string) error
}

type client struct {
	url        url.URL
	username   string
	password   string
	httpClient *http.Client
}

func NewHTTPClient(conf *HTTPConfig) (Client, error) {
	u, err := url.Parse(conf.Addr)
	if err != nil {
		return nil, err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		m := fmt.Sprintf("Unsupported protocol scheme: %s, your address"+
			" must start with http:// or https://", u.Scheme)
		return nil, errors.New(m)
	}

	return &client{
		url:      *u,
		username: conf.Username,
		password: conf.Password,
		httpClient: &http.Client{
			Timeout: conf.Timeout,
		},
	}, nil
}

func (c *client) Ping(timeout time.Duration) (time.Duration, string, error) {
	return 0, "", nil
}

func setReqHeader(req *http.Request) *http.Request {
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	return req
}

func (c *client) Login() error {
	req, err := http.NewRequest("POST", c.url.String()+"/login", nil)
	if err != nil {
		fmt.Printf("failed to post request", err)
	}

	req = setReqHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("failed to login: \n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *client) Logout() error {
	return nil
}

func (c *client) Submit(path string) error {
	return nil
}

func (c *client) Run(path string) error {
	u := c.url
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req = setReqHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

func (c *client) Search(v interface{}) error {
	return nil
}
