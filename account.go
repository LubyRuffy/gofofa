package gofofa

type DeductMode int // 扣费模式
const (
	DeductModeFree  DeductMode = 0 // 不自动扣费
	DeductModeFCoin            = 1 // 自动扣除F币
)

func ParseDeductMode(v string) DeductMode {
	switch v {
	case "0", "DeductModeFree":
		return DeductModeFree
	case "1", "DeductModeFCoin":
		return DeductModeFCoin
	default:
		panic("unknown deduct mode")
	}
}

// AccountInfo fofa account info
type AccountInfo struct {
	Error    bool   `json:"error"`     // error or not
	ErrMsg   string `json:"errmsg"`    // error string message
	FCoin    int    `json:"fcoin"`     // fcoin count
	IsVIP    bool   `json:"isvip"`     // is vip
	VIPLevel int    `json:"vip_level"` // vip level
}

// AccountInfo fetch account info from fofa
func (c *Client) AccountInfo() (ac AccountInfo, err error) {
	err = c.fetch("info/my", nil, &ac)
	return
}

// freeSize 获取可以免费使用的数据量
func (c *Client) freeSize() int {
	if !c.Account.IsVIP {
		return 0
	}

	switch c.Account.VIPLevel {
	case 0:
		return 0
	case 1:
		return 100
	case 2:
		return 10000
	case 3:
		return 100000
	}
	return 0
}
