package gofofa

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
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

func readAll(reader io.Reader, size int) ([]byte, error) {
	if size <= 0 {
		size = int(math.Max(float64(size), 65535))
	}
	buffer := bytes.NewBuffer(make([]byte, 0, size))
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// just fetch fofa body, no need to unmarshal
func (c *Client) fetchBody(apiURI string, params map[string]string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	fullURL := c.buildURL(apiURI, params)
	c.logger.Debugf("fetch fofa: %s", apiURI)
	//c.logger.Debugf("fetch fofa: %s", fullURL)

	req, err = http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Accept-Encoding", "gzip")
	//requestDump, _ := httputil.DumpRequestOut(req, false)
	//logrus.Debugln(string(requestDump))

	resp, err = c.httpClient.Do(req)
	if err != nil {
		if !c.accountDebug {
			// 替换账号明文信息
			if e, ok := err.(*url.Error); ok {
				newClient := c
				newClient.Email = "<email>"
				newClient.Key = "<key>"
				e.URL = newClient.buildURL(apiURI, params)
				err = e
			}
		}
		return
	}
	defer resp.Body.Close()

	contentLength := 0
	if v := resp.Header.Get("Content-Length"); len(v) > 0 {
		// 这个地方不可能出错，因为在httpclient get过程中进行了合法性校验
		contentLength, _ = strconv.Atoi(v)
	}
	encoding := resp.Header.Get("Content-Encoding")
	// 取body
	body, err = readAll(resp.Body, contentLength)
	if err != nil {
		return
	}

	switch encoding {
	case "gzip":
		var reader *gzip.Reader
		reader, err = gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return
		}
		body, _ = readAll(reader, 0) // 这里的err也是不可能的，不用关注
		//case "deflate":
		//	reader1 := flate.NewReader(bytes.NewReader(body))
		//	body, err = readAll(reader1, 0)
		//	if err != nil {
		//		return
		//	}
		//	reader1.Close()
	}
	if encoding == "gzip" {

	}

	//respDump, _ := httputil.DumpResponse(resp, false)
	//logrus.Debugln(string(respDump))

	return
}

// Fetch http request and parse as json return to v
func (c *Client) Fetch(apiURI string, params map[string]string, v interface{}) (err error) {
	content, err := c.fetchBody(apiURI, params)
	if err != nil {
		return
	}

	if err = json.Unmarshal(content, v); err != nil {
		return
	}
	return
}
