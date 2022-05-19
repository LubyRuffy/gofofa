package gofofa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Fetch(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:55")
	assert.Error(t, err)

	ts := httptest.NewServer(http.HandlerFunc(queryHander))
	defer ts.Close()

	account := validAccounts[1]
	fofaURL := ts.URL + "/?email=" + account.Email + "&key=" + account.Key + "&version=v1"
	cli, err := NewClient(fofaURL)

	// 解析异常
	var a map[string]interface{}
	err = cli.Fetch("/", nil, &a)
	assert.Error(t, err)
}
