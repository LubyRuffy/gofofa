package gofofa

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidQuery = errors.New("query is not valid")
)

// StatsResults /search/stats api results
type StatsResults struct {
	Error          bool                    `json:"error"`
	Errmsg         string                  `json:"errmsg"`
	Distinct       map[string]interface{}  `json:"distinct"`
	Aggs           map[string][]*StatsItem `json:"aggs"`
	LastUpdateTime string                  `json:"lastupdatetime"`
}

type CertObject struct {
	CN            string   `json:"cn"`
	Organizations []string `json:"org"`
}

type CertDetail struct {
	Subject     CertObject  `json:"subject"`
	Issuer      *CertObject `json:"issuer"`
	RootDomains []string    `json:"domain"`
	IsExpired   bool        `json:"is_expired"`
	IsValid     bool        `json:"is_valid"`
	NotAfter    string      `json:"not_after"`
	NotBefore   string      `json:"not_before"`
}

type IPDetail struct {
	Domains []string `json:"domains"`
}

type ASNDetail struct {
	Org *string `json:"org"`
}
type IconDetail struct {
	IconBase64 *string `json:"icon_base64"`
}

type Detail struct {
	CertDetail
	IPDetail
	ASNDetail
	IconDetail
}

// StatsItem one stats item
type StatsItem struct {
	Name   string         `json:"name"`
	Count  int            `json:"count"`
	Uniq   map[string]int `json:"uniq,omitempty"`
	Detail *Detail        `json:"detail,omitempty"`
}

// StatsObject one stats object
type StatsObject struct {
	Name  string
	Items []*StatsItem
}

// Stats aggs fofa host data
// query fofa query string
// size data size
// fields' field of fofa host struct
func (c *Client) Stats(query string, size int, fields []string) (res []StatsObject, err error) {
	if len(fields) == 0 {
		fields = []string{"title", "country"}
	}

	if query == "" {
		return nil, ErrInvalidQuery
	}

	var sr StatsResults
	err = c.Fetch("search/stats",
		map[string]string{
			"qbase64":    base64.StdEncoding.EncodeToString([]byte(query)),
			"size":       strconv.Itoa(size),
			"fields":     strings.Join(fields, ","),
			"full":       "false", // 是否全部数据，非一年内
			"uniq_count": "ip",    // uniq 统计的字段，默认为ip
			"detail":     "true",  // 是否显示详情？目前仅仅在cert.sn的场景下有效
		},
		&sr)
	if err != nil {
		return
	}
	if len(sr.Errmsg) > 0 {
		err = errors.New(sr.Errmsg)
		return
	}

	for _, rawField := range fields {
		field := rawField
		// 有些字段要进行改名
		switch rawField {
		case "country":
			field = "countries"
		}

		if v, ok := sr.Aggs[field]; ok {
			so := StatsObject{
				Name:  rawField,
				Items: v,
			}
			res = append(res, so)
		}
	}

	return
}
