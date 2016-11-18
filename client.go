package goketo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

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
	authLock sync.Mutex
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
	client := &Client{
		client: &http.Client{
			Transport: &bearerRoundTripper{
				clientID:     clientID,
				clientSecret: ClientSecret,
			},
		},
		endpoint: endpoint,
		identity: identity,
		version:  version,
	}

	if err := client.RefreshToken(); err != nil {
		return nil, err
	}
	return client, nil
}

// RefreshToken refreshes the auth token provided by the Marketo API.
func (c *Client) RefreshToken() error {
	// Make request for token
	resp, err := c.client.Get(c.identity + "oauth/token?grant_type=client_credentials")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("an error occured while fetching data: %v", resp)
	}

	var auth AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		logrus.Errorf("Could not convert response: %v", err)
		return err
	}
	c.authLock.Lock()
	defer c.authLock.Unlock()
	c.auth = &auth
	return nil
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
	c.authLock.Lock()
	logrus.Debug("Token: ", c.auth.Token)
	req.Header.Add("Authorization", "Bearer "+c.auth.Token)
	c.authLock.Unlock()
	return c.do(req)
}

// Post to resource string the data provided
func (c *Client) Post(resource string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.endpoint+resource, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	c.authLock.Lock()
	logrus.Debug("Token: ", c.auth.Token)
	req.Header.Add("Authorization", "Bearer "+c.auth.Token)
	c.authLock.Unlock()
	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}
