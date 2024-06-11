package github

import "net/http"

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) {
		if token != "" {
			c.token = token
		}
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.client = client
		}
	}
}

func WithSilent(silent bool) Option {
	return func(c *Client) {
		c.silent = silent
	}
}
