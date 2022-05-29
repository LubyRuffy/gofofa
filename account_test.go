package gofofa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestParseDeductMode(t *testing.T) {
	assertPanic(t, func() {
		ParseDeductMode("abc")
	})

	assert.Equal(t, DeductModeFree, ParseDeductMode("DeductModeFree"))
	assert.Equal(t, DeductModeFree, ParseDeductMode("0"))
	assert.Equal(t, DeductModeFCoin, ParseDeductMode("DeductModeFCoin"))
	assert.Equal(t, DeductModeFCoin, ParseDeductMode("1"))
}

func TestAccountInfo_String(t *testing.T) {
	ai := AccountInfo{
		Error:    false,
		IsVIP:    true,
		VIPLevel: 3,
		FCoin:    0,
	}
	assert.Equal(t, `{
  "error": false,
  "fcoin": 0,
  "isvip": true,
  "vip_level": 3
}`, ai.String())
}

func TestClient_AccountInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(queryHander))
	defer ts.Close()

	var cli *Client
	var err error
	_, err = NewClient(WithURL(ts.URL + "?email=a@a.com&key=wrong"))
	assert.Contains(t, err.Error(), "[-700] Account Invalid")

	// 注册用户
	var account accountInfo
	account = validAccounts[0]
	cli, err = NewClient(WithURL(ts.URL + "?email=" + account.Email + "&key=" + account.Key))
	assert.Nil(t, err)
	assert.False(t, cli.Account.IsVIP)
	assert.Equal(t, 0, cli.Account.FCoin)
	assert.Equal(t, 0, cli.freeSize())

	// 普通会员
	account = validAccounts[1]
	cli, err = NewClient(WithURL(ts.URL + "?email=" + account.Email + "&key=" + account.Key))
	assert.Nil(t, err)
	assert.True(t, cli.Account.IsVIP)
	assert.Equal(t, 1, cli.Account.VIPLevel)
	assert.Equal(t, 10, cli.Account.FCoin)
	assert.Equal(t, 100, cli.freeSize())

	// 高级会员
	account = validAccounts[2]
	cli, err = NewClient(WithURL(ts.URL + "?email=" + account.Email + "&key=" + account.Key))
	assert.Nil(t, err)
	assert.True(t, cli.Account.IsVIP)
	assert.Equal(t, 2, cli.Account.VIPLevel)
	assert.Equal(t, 0, cli.Account.FCoin)
	assert.Equal(t, 10000, cli.freeSize())

	// 企业会员
	account = validAccounts[3]
	cli, err = NewClient(WithURL(ts.URL + "?email=" + account.Email + "&key=" + account.Key))
	assert.Nil(t, err)
	assert.True(t, cli.Account.IsVIP)
	assert.Equal(t, 3, cli.Account.VIPLevel)
	assert.Equal(t, 0, cli.Account.FCoin)
	assert.Equal(t, 100000, cli.freeSize())
}
