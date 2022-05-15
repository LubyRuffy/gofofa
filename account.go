package gofofa

// AccountInfo fofa account info
type AccountInfo struct {
	Error  bool   `json:"error"`  // error or not
	ErrMsg string `json:"errmsg"` // error string message
	FCoin  int    `json:"fcoin"`  // error string message
}

// AccountInfo fetch account info from fofa
func (c *Client) AccountInfo() (ac AccountInfo, err error) {
	err = c.fetch("info/my", &ac)
	return
}
