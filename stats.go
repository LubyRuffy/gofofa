package gofofa

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// StatsResults /search/stats api results
type StatsResults struct {
	Error          bool                   `json:"error"`
	Errmsg         string                 `json:"errmsg"`
	TraceId        string                 `json:"trace_id"`
	Distinct       map[string]interface{} `json:"distinct"`
	Aggs           map[string]interface{} `json:"aggs"`
	LastUpdateTime string                 `json:"lastupdatetime"`
}

func (s StatsResults) SetTraceId(traceId string) {
	s.TraceId = traceId
}

// StatsItem one stats item
type StatsItem struct {
	Name  string
	Count int
}

// StatsObject one stats object
type StatsObject struct {
	Name  string
	Items []StatsItem
}

// Stats aggs fofa host data
// query fofa query string
// size data size
// fields' field of fofa host struct
func (c *Client) Stats(query string, size int, fields []string) (res []StatsObject, err error) {
	if len(fields) == 0 {
		fields = []string{"title", "country"}
	}

	var sr StatsResults
	err = c.Fetch("search/stats",
		map[string]string{
			"qbase64": base64.StdEncoding.EncodeToString([]byte(query)),
			"size":    strconv.Itoa(size),
			"fields":  strings.Join(fields, ","),
			"full":    "false", // 是否全部数据，非一年内
		},
		&sr)
	if err != nil {
		if c.traceId {
			err = fmt.Errorf("[%s]%s", sr.TraceId, err.Error())
		}
		return
	}

	// 报错，退出
	if len(sr.Errmsg) > 0 {
		if c.traceId {
			err = errors.New(sr.Errmsg + " trace id: " + sr.TraceId)
		} else {
			err = errors.New(sr.Errmsg)
		}
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
			if objArray, ok := v.([]interface{}); ok && len(objArray) > 0 {
				so := StatsObject{
					Name: rawField,
				}
				for _, obj := range objArray {
					obj := obj.(map[string]interface{})
					so.Items = append(so.Items, StatsItem{
						Name:  obj["name"].(string),
						Count: int(obj["count"].(float64)),
					})
				}
				res = append(res, so)
			}
		}
	}

	return
}
