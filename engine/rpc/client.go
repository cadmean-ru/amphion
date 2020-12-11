package rpc

import "net/http"

type Client struct {
	url        string
	httpClient *http.Client
}

var instance *Client

func Initialize(url string) {
	if instance != nil {
		return
	}

	instance = &Client{
		url: url,
		httpClient: &http.Client{},
	}
}

func F(fName string) *CallBuilder {
	return newCallBuilder(fName)
}