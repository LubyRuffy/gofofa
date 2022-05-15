package gofofa

import (
	"encoding/base64"
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
)

type HostResults struct {
	Mode    string
	Error   bool
	Errmsg  string
	Query   string
	Page    int
	Size    int
	Results [][]string
}

// HostSearch search fofa host data
// query fofa query string
// size data size: -1 means all，0 means just data total info, >0 means actual size
// fields field of fofa host struct
func (c *Client) HostSearch(query string, size int, fields []string) (res [][]string, err error) {
	// check level
	if c.freeSize() == 0 {
		// 不是会员
		if c.Account.FCoin < 1 {
			return nil, errors.New("insufficient privileges") // 等级不够，fcoin也不够
		}
		if c.DeductMode != DeductModeFCoin {
			return nil, errors.New("insufficient privileges, try to set mode to 1(DeductModeFCoin)") // 等级不够，fcoin也不够
		}
	} else if size > c.freeSize() {
		// 是会员，但是取的数量比免费的大
		switch c.DeductMode {
		case DeductModeFree:
			size = c.freeSize()
			log.Println("[WARNING] size is larger than your account free limit, ",
				"just fetch %d instead, if you want deduct fcoin automatically, set mode to 1(DeductModeFCoin) manually")
		}
	}

	page := 1
	perPage := int(math.Min(float64(size), 1000)) // 最多一次取1000
	if len(fields) == 0 {
		fields = []string{"ip", "port"}
	}

	// 分页取数据
	for {
		var hr HostResults
		err = c.fetch("search/all",
			map[string]string{
				"qbase64": base64.StdEncoding.EncodeToString([]byte(query)),
				"size":    strconv.Itoa(perPage),
				"page":    strconv.Itoa(page),
				"fields":  strings.Join(fields, ","),
				"full":    "false", // 是否全部数据，非一年内
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

		// 无数据
		if len(hr.Results) == 0 {
			break
		}

		res = append(res, hr.Results...)

		// 数据填满了，完成
		if size <= len(res) {
			break
		}

		// 数据已经没有了
		if len(hr.Results) < perPage {
			break
		}
	}

	return
}
