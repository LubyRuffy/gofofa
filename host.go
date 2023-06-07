package gofofa

import (
	"context"
	"encoding/base64"
	"errors"
	"math"
	"strconv"
	"strings"
)

// HostResults /search/all api results
type HostResults struct {
	Mode    string      `json:"mode"`
	Error   bool        `json:"error"`
	Errmsg  string      `json:"errmsg"`
	Query   string      `json:"query"`
	Page    int         `json:"page"`
	Size    int         `json:"size"` // 总数
	Results interface{} `json:"results"`
}

// HostStatsData /host api results
type HostStatsData struct {
	Error       bool     `json:"error"`
	Errmsg      string   `json:"errmsg"`
	Host        string   `json:"host"`
	IP          string   `json:"ip"`
	ASN         int      `json:"asn"`
	ORG         string   `json:"org"`
	Country     string   `json:"country_name"`
	CountryCode string   `json:"country_code"`
	Protocols   []string `json:"protocol"`
	Ports       []int    `json:"port"`
	Categories  []string `json:"category"`
	Products    []string `json:"product"`
	UpdateTime  string   `json:"update_time"`
}

// SearchOptions options of search, for post processors
type SearchOptions struct {
	FixUrl    bool   // each host fix as url, like 1.1.1.1,80 will change to http://1.1.1.1, https://1.1.1.1:8443 will no change
	UrlPrefix string // default is http://
}

// SecondaryOptions options used while searching, use default option for empty input
type SecondaryOptions struct {
	Full bool // search result for over a year
}

// fixHostToUrl 替换host为url
func fixHostToUrl(res [][]string, fields []string, hostIndex int, urlPrefix string) [][]string {
	newRes := make([][]string, 0, len(res))
	for _, row := range res {
		newRow := make([]string, 0, len(fields))
		for j, r := range row {
			if j == hostIndex {
				if !strings.Contains(r, "://") {
					r = urlPrefix + r
				}
			}
			newRow = append(newRow, r)
		}
		newRes = append(newRes, newRow)
	}
	return newRes
}

// HostSearch search fofa host data
// query fofa query string
// size data size: -1 means all，0 means just data total info, >0 means actual size
// fields of fofa host search
// sOptions - Full: fetch result over a year, default false
func (c *Client) HostSearch(query string, size int, fields []string, sOptions *SecondaryOptions, options ...SearchOptions) (res [][]string, err error) {
	// assign default value for secondary options
	if sOptions == nil {
		sOptions = &SecondaryOptions{Full: false}
	}

	freeSize := c.freeSize()
	// check level
	if freeSize == 0 {
		// 不是会员
		if c.Account.FCoin < 1 {
			return nil, errors.New("insufficient privileges") // 等级不够，fcoin也不够
		}
		if c.DeductMode != DeductModeFCoin {
			return nil, errors.New("insufficient privileges, try to set mode to 1(DeductModeFCoin)") // 等级不够，fcoin也不够
		}
	} else if freeSize == -1 {
		// unknown vip level, skip mode check
	} else if size > c.freeSize() {
		// 是会员，但是取的数量比免费的大
		switch c.DeductMode {
		case DeductModeFree:
			// 防止 freesize = -1，取 size 和 freesize 的最大值
			if freeSize <= 0 {
				size = int(math.Max(float64(freeSize), float64(size)))
			} else {
				size = freeSize
			}
			c.logger.Warnf("size is larger than your account free limit, "+
				"just fetch %d instead, if you want deduct fcoin automatically, set mode to 1(DeductModeFCoin) manually", size)
		}
	}

	page := 1
	perPage := int(math.Min(float64(size), 1000)) // 最多一次取1000

	// 一次取所有数据，perPage 默认给 1000
	if size == -1 {
		perPage = 1000
	}

	if len(fields) == 0 {
		fields = []string{"ip", "port"}
	}

	// 分页取数据
	for {
		if ctx := c.GetContext(); ctx != nil {
			// 确认是否需要退出
			select {
			case <-c.GetContext().Done():
				err = context.Canceled
				return
			default:
			}
		}

		var hr HostResults
		err = c.Fetch("search/all",
			map[string]string{
				"qbase64": base64.StdEncoding.EncodeToString([]byte(query)),
				"size":    strconv.Itoa(perPage),
				"page":    strconv.Itoa(page),
				"fields":  strings.Join(fields, ","),
				"full":    strconv.FormatBool(sOptions.Full), // 是否全部数据，非一年内
			},
			&hr)
		if err != nil {
			return
		}

		// 报错，退出
		if len(hr.Errmsg) > 0 {
			err = errors.New(hr.Errmsg)
			break
		}

		var results [][]string
		if v, ok := hr.Results.([]interface{}); ok {
			// 无数据
			if len(v) == 0 {
				break
			}
			for _, result := range v {
				if vStrSlice, ok := result.([]interface{}); ok {
					var newSlice []string
					for _, vStr := range vStrSlice {
						newSlice = append(newSlice, vStr.(string))
					}
					results = append(results, newSlice)
				} else if vStr, ok := result.(string); ok {
					results = append(results, []string{vStr})
				}
			}
		} else {
			break
		}

		if c.OnResults != nil {
			c.OnResults(results)
		}

		res = append(res, results...)

		// 数据填满了，完成
		if size != -1 && size <= len(res) {
			break
		}

		// 数据已经没有了
		if len(results) < perPage {
			break
		}

		page++ // 翻页
	}

	// 后处理
	if len(options) > 0 && options[0].FixUrl {
		urlPrefix := options[0].UrlPrefix
		if urlPrefix == "" {
			urlPrefix = "http://"
		}
		for i, f := range fields {
			if f == "host" {
				res = fixHostToUrl(res, fields, i, urlPrefix)
				break
			}
		}
	}

	return
}

// HostSize fetch query matched host count
func (c *Client) HostSize(query string) (count int, err error) {
	var hr HostResults
	err = c.Fetch("search/all",
		map[string]string{
			"qbase64": base64.StdEncoding.EncodeToString([]byte(query)),
			"size":    "1",
			"page":    "1",
			"full":    "false", // 是否全部数据，非一年内
		},
		&hr)
	if err != nil {
		return
	}
	count = hr.Size
	return
}

// HostStats fetch query matched host count
func (c *Client) HostStats(host string) (data HostStatsData, err error) {
	err = c.Fetch("host/"+host, nil, &data)
	if err != nil {
		return
	}
	return
}
