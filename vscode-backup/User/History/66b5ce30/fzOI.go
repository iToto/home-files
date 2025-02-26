package client

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//go:generate mockgen --build_flags=-mod=mod -package=client -destination=mock_httpclient.go . Client,Doer,Middleware

// Doer interface has the method required to use a type as custom http client.
// The net/*http.Client type satisfies this interface.
// Note that the caller is responsible for closing the response body
type Doer interface {
	// Do makes an HTTP request and returns an *http.Response or an errors
	Do(*http.Request) (*http.Response, error)
}

// Client Is a generic HTTP client interface
type Client interface {
	Doer
	// DoJSON makes an HTTP request and decodes the json body into a given struct
	DoJSON(req *http.Request, respObj interface{}) error
}

// Middleware is a generic interface for http middleware
type Middleware interface {
	OnRequestStart(*http.Request)
	OnRequestEnd(*http.Request, *http.Response)
	OnError(*http.Request, error)
}

const (
	dialTimeout = time.Second * 10
	tlsTimeout  = time.Second * 5
)

// HTTPClient is an http client that supports json, it's an implementation of the Client interface
type HTTPClient struct {
	client      Doer
	middlewares []Middleware
}

// NewHTTPClient returns a new instance of http Client
func NewHTTPClient(opts ...Option) *HTTPClient {
	var c HTTPClient

	for _, opt := range opts {
		opt(&c)
	}

	if c.client == nil {
		c.client = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   dialTimeout,
					KeepAlive: dialTimeout,
				}).Dial,
				TLSHandshakeTimeout: tlsTimeout,
			},
		}
	}

	return &c
}

// DoJSON makes an HTTP request and decodes the json body into a given struct
func (c *HTTPClient) DoJSON(req *http.Request, respObj interface{}) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&respObj)
}

// Do makes an HTTP request and returns an *http.Response or an errors
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	c.reportRequestStart(req)
	resp, err = c.client.Do(req)
	if err != nil {
		c.reportError(req, err)
		return nil, err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		c.reportError(req, err)
		respErr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}
		return nil, &HTTPError{URL: req.URL.String(), StatusCode: resp.StatusCode, Body: string(respErr)}
	}
	c.reportRequestEnd(req, resp)
	return resp, err
}

func (c *HTTPClient) reportRequestStart(request *http.Request) {
	for _, plugin := range c.middlewares {
		plugin.OnRequestStart(request)
	}
}

func (c *HTTPClient) reportError(request *http.Request, err error) {
	for _, plugin := range c.middlewares {
		plugin.OnError(request, err)
	}
}

func (c *HTTPClient) reportRequestEnd(request *http.Request, response *http.Response) {
	for _, plugin := range c.middlewares {
		plugin.OnRequestEnd(request, response)
	}
}
