package pipeparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	// 单个
	assert.Equal(t, "a()\n", NewParser().Parse(`a()`))

	// 两个
	assert.Equal(t, "a()\nb()\n", NewParser().Parse(`a() | b()`))

	// 多个
	assert.Equal(t, "a()\nb()\nc()\nd()\ne()\n", NewParser().Parse("a( )| b()|  c  ()  |     d(   )    |e(\n\n)"))

	// 换行
	assert.Equal(t, "a()\nb()\n", NewParser().Parse("a () \n| b()"))

	// 参数
	assert.Equal(t, "a()\n", NewParser().Parse(`a(1)`))
	assert.Equal(t, "a()\n", NewParser().Parse(`a( 1   ,   2 )`))

	// 嵌套
	// assert.Equal(t, "a(b())\n", NewParser().Parse(`a(b(1))`))
}
