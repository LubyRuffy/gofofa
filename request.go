package gofofa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// params is key=>value for query, auto encoded with uri escape
func (c *Client) buildURL(apiURI string, params map[string]string) string {
	fullURL := fmt.Sprintf("%s/api/%s/%s?", c.Server, c.APIVersion, apiURI)
	ps := url.Values{}
	ps.Set("email", c.Email)
	ps.Set("key", c.Key)
	for k, v := range params {
		ps.Set(k, v)
	}
	return fullURL + ps.Encode()
}

// http request and parse as json return to v
func (c *Client) fetch(apiURI string, params map[string]string, v interface{}) (err error) {
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest("GET", c.buildURL(apiURI, params), nil)
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(content, v); err != nil {
		return
	}
	return
}
