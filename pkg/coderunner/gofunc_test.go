package coderunner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoFunction_Register(t *testing.T) {
	var err error
	gf := &GoFunction{}

	// 正确
	err = gf.Register("abc", func() {})
	assert.Nil(t, err)

	// 函数名称格式异常
	err = gf.Register("111", func() {})
	assert.Error(t, err)
	err = gf.Register("abc()", func() {})
	assert.Error(t, err)

	// 函数定义格式异常
	err = gf.Register("abc", "123")
	assert.Error(t, err)

}
