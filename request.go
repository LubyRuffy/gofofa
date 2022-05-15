package gofofa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) buildURL(apiURI string) string {
	return fmt.Sprintf("%s/api/%s/%s?email=%s&key=%s", c.Server, c.APIVersion, apiURI, c.Email, c.Key)
}

// http request and parse as json return to v
func (c *Client) fetch(apiURI string, v interface{}) (err error) {
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest("GET", c.buildURL(apiURI), nil)
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
