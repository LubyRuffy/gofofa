package gofofa

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Stats(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(queryHander))
	defer ts.Close()

	var cli *Client
	var err error
	var account accountInfo
	var res []StatsObject

	account = validAccounts[1]
	cli, err = NewClient(WithURL(ts.URL + "?email=" + account.Email + "&key=" + account.Key))
	assert.Nil(t, err)

	// 错误
	res, err = cli.Stats("port=80", 0, []string{"title"})
	assert.Error(t, err)

	// 正确
	res, err = cli.Stats("port=80", 5, []string{"title"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 5, len(res[0].Items))
	assert.Equal(t, "title", res[0].Name)
	assert.Equal(t, 25983408, res[0].Items[0].Count)

	// 默认字段
	res, err = cli.Stats("port=80", 5, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, 5, len(res[0].Items))
	assert.Equal(t, "title", res[0].Name)
	assert.Equal(t, "301 Moved Permanently", res[0].Items[0].Name)
	assert.Equal(t, 25983454, res[0].Items[0].Count)
	assert.Equal(t, "country", res[1].Name)
	assert.Equal(t, "United States of America", res[1].Items[0].Name)
	assert.Equal(t, 154746752, res[1].Items[0].Count)

	// 测试fields参数为空的情况
	res, err = cli.Stats("port=80", 5, []string{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "title", res[0].Name)
	assert.Equal(t, "country", res[1].Name)

	// 测试fields参数包含多个字段的情况
	res, err = cli.Stats("port=80", 5, []string{"title", "country"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "title", res[0].Name)
	assert.Equal(t, "country", res[1].Name)

	// 测试size参数为0的情况
	res, err = cli.Stats("port=80", 0, []string{"title"})
	assert.Error(t, err)

	// 测试size参数为负数的情况
	res, err = cli.Stats("port=80", -1, []string{"title"})
	assert.NoError(t, err)

	// 测试query参数为空的情况
	res, err = cli.Stats("", 5, []string{"title"})
	assert.EqualError(t, err, ErrInvalidQuery.Error())

	// 测试CertDetail字段的解析
	res, err = cli.Stats(`port=80`, 5, []string{"cert.sn"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "cert.sn", res[0].Name)
	assert.Equal(t, 5, len(res[0].Items))
	assert.Equal(t, "fofa.info", res[0].Items[0].Detail.RootDomains[0])

	// 请求失败
	cli = &Client{
		Server:     "http://fofa.info:66666",
		httpClient: &http.Client{},
		logger:     logrus.New(),
	}
	res, err = cli.Stats("port=80", 5, nil)
	assert.Error(t, err)
}
