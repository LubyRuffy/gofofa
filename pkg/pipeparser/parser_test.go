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
	assert.Equal(t, "a(1)\n", NewParser().Parse(`a(1)`))
	assert.Equal(t, "a(1, 2)\n", NewParser().Parse(`a( 1   ,   2 )`))

	// 嵌套
	assert.Equal(t, "a(b(1))\n", NewParser().Parse(`a(b(1))`))

	// 字符串参数
	assert.Equal(t, "a(\"abc\")\n", NewParser().Parse(`a("abc")`))
	assert.Equal(t, `a("abc\"123")`+"\n", NewParser().Parse(`a("abc\"123")`))

	// 0x
	assert.Equal(t, "a(0x10)\n", NewParser().Parse(`a(0x10)`))

	// char
	assert.Equal(t, "a('a')\n", NewParser().Parse(`a('a')`))
}
