package goketo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Requester is the interface for all client calls - This allows easy test 'mocks'
type Requester interface {
	do(*http.Request) ([]byte, error)
	Get(string) ([]byte, error)
	Post(string, []byte) ([]byte, error)
}

// Client http client tracker
type Client struct {
	client   *http.Client
	endpoint string
	identity string
	version  string
	auth     *AuthToken
}

// bearerRoundTripper wrapper for query params
type bearerRoundTripper struct {
	Delegate     http.RoundTripper
	clientID     string
	clientSecret string
}

// AuthToken holds data from Auth request
type AuthToken struct {
	Token   string `json:"access_token"`
	Type    string `json:"token_type"`
	Expires int    `json:"expires_in"` // in seconds
	Scope   string `json:"scope"`
}

func (b *bearerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if b.Delegate == nil {
		b.Delegate = http.DefaultTransport
	}
	values := req.URL.Query()
	values.Add("client_id", b.clientID)
	values.Add("client_secret", b.clientSecret)
	req.URL.RawQuery = values.Encode()
	return b.Delegate.RoundTrip(req)
}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

// NewAuthClient request application/json
func NewAuthClient(clientID string, ClientSecret string, ClientEndpoint string) (*Client, error) {
	// Endpoint: /identity/oauth/token?grant_type=client_credentials
	version := "v1"
	var endpoint string
	var identity string

	// Check if endpoint has proper protocol
	if strings.HasPrefix(ClientEndpoint, "http") {
		endpoint = ClientEndpoint + "/rest/" + version + "/"
		identity = ClientEndpoint + "/identity/"
	} else {
		endpoint = "https://" + ClientEndpoint + "/rest/" + version + "/"
		identity = "https://" + ClientEndpoint + "/identity/"
	}
	// Add credentials to the request
	client := &http.Client{
		Transport: &bearerRoundTripper{
			clientID:     clientID,
			clientSecret: ClientSecret,
		},
	}
	// Make request for token
	resp, err := client.Get(identity + "oauth/token?grant_type=client_credentials")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var auth AuthToken
	if resp.StatusCode == 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Could not convert response: %v", err)
		}
		err = json.Unmarshal(data, &auth)
	} else {
		logrus.Errorf("An error occured while fetching data: %v", resp)
	}

	return &Client{
		client:   client,
		endpoint: endpoint,
		identity: identity,
		version:  version,
		auth:     &auth,
	}, nil
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Received unexpected status %d while trying to retrieve the server data with \"%s\"", resp.StatusCode, string(body))
		return nil, err
	}
	return body, nil
}

// Get resource string
func (c *Client) Get(resource string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.endpoint+resource, nil)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Token: ", c.auth.Token)
	req.Header.Add("Authorization", "Bearer "+c.auth.Token)
	return c.do(req)
}

// Post to resource string the data provided
func (c *Client) Post(resource string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.endpoint+resource, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	logrus.Debug("Token: ", c.auth.Token)
	req.Header.Add("Authorization", "Bearer "+c.auth.Token)
	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}
