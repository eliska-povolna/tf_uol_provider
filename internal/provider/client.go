package provider

import (
	"encoding/base64"
	"net/http"
)

// Client structure to store API credentials
type Client struct {
	Email      string
	Token      string
	HttpClient *http.Client
}

func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(c.Email + ":" + c.Token))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return c.HttpClient.Do(req)
}
