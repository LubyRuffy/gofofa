/*Package gofofa fofa client in Go

env settings:
- FOFA_CLIENT_URL full fofa connnection string, format: <url>/?email=<email>&key=<key>&version=<v2>
- FOFA_SERVER fofa server
- FOFA_EMAIL fofa account email
- FOFA_KEY fofa account key
*/
package gofofa

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

const (
	defaultServer     = "https://fofa.info"
	defaultAPIVersion = "v1"
)

// Client of fofa connection
type Client struct {
	Server     string // can set local server for debugging, format: <scheme>://<host>
	APIVersion string // api version
	Email      string // fofa email
	Key        string // fofa key

	Account    AccountInfo // fofa account info
	DeductMode DeductMode  // 扣费提醒默认

	httpClient *http.Client //
	logger     *logrus.Logger
}

// Update merge config from config url
func (c *Client) Update(configURL string) error {
	u, err := url.Parse(configURL)
	if err != nil {
		return err
	}

	c.Server = u.Scheme + "://" + u.Host
	if u.Query().Has("email") {
		c.Email = u.Query().Get("email")
	}

	if u.Query().Has("key") {
		c.Key = u.Query().Get("key")
	}

	if u.Query().Has("version") {
		c.APIVersion = u.Query().Get("version")
	}

	return nil
}

// URL generate fofa connection url string
func (c *Client) URL() string {
	return fmt.Sprintf("%s/?email=%s&key=%s&version=%s", c.Server, c.Email, c.Key, c.APIVersion)
}

type ClientOption func(c *Client) error

// WithURL configURL format: <url>/?email=<email>&key=<key>&version=<v2>&tlsdisabled=false&debuglevel=0
func WithURL(configURL string) ClientOption {
	return func(c *Client) error {
		// merge from config
		if len(configURL) > 0 {
			return c.Update(configURL)
		}
		return nil
	}
}

// WithLogger set logger
func WithLogger(logger *logrus.Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// NewClient from fofa connection string to config
// and with env config merge
func NewClient(options ...ClientOption) (*Client, error) {
	// read from env
	c, err := newClientFromEnv()
	if err != nil {
		return c, err
	}

	c.logger = logrus.New()
	for _, opt := range options {
		err = opt(c)
		if err != nil {
			return c, err
		}
	}

	// fetch one time to make sure network is ok
	c.httpClient = &http.Client{}
	c.Account, err = c.AccountInfo()
	if err != nil {
		return c, err
	}

	if c.Account.Error {
		return c, fmt.Errorf("auth failed: '%s', make sure key is valid", c.Account.ErrMsg)
	}

	return c, nil
}
