package gofofa

import (
	"github.com/stretchr/testify/assert"
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
