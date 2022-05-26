package workflowast

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
)

func TestNewParser(t *testing.T) {
	// 单个
	assert.Equal(t, "a()\n", NewParser().MustParse(`a()`))

	// 两个
	assert.Equal(t, "a()\nb()\n", NewParser().MustParse(`a() | b()`))

	// 多个
	assert.Equal(t, "a()\nb()\nc()\nd()\ne()\n", NewParser().MustParse("a( )| b()|  c  ()  |     d(   )    |e(\n\n)"))

	// 换行
	assert.Equal(t, "a()\nb()\n", NewParser().MustParse("a () \n| b()"))

	// 参数
	assert.Equal(t, "a(1)\n", NewParser().MustParse(`a(1)`))
	assert.Equal(t, "a(1, 2)\n", NewParser().MustParse(`a( 1   ,   2 )`))

	// 嵌套
	assert.Equal(t, "a(b(1))\n", NewParser().MustParse(`a(b(1))`))

	// 字符串参数
	assert.Equal(t, "a(\"abc\")\n", NewParser().MustParse(`a("abc")`))
	assert.Equal(t, "a(`abc`)\n", NewParser().MustParse("a(`abc`)"))
	assert.Equal(t, `a("abc\"123")`+"\n", NewParser().MustParse(`a("abc\"123")`))

	// 0x
	assert.Equal(t, "a(0x10)\n", NewParser().MustParse(`a(0x10)`))

	// char
	assert.Equal(t, "a('a')\n", NewParser().MustParse(`a('a')`))

}

func TestParser_Fork(t *testing.T) {
	// 测试分叉fork的格式
	assert.Equal(t, "Fork(\"cut(`ip`)\")\nFork(\"cut(`port`)\")\n",
		NewParser().MustParse("[ cut(`ip`) & cut(`port`) ]"))
	assert.Equal(t, "Fork(\"cut(\\\"ip\\\")\")\nFork(\"cut(`port`)\")\n",
		NewParser().MustParse("[ cut(\"ip\") & cut(`port`) ]"))

	// 多层嵌套
	assert.Equal(t, "Fork(\"cut(`ip`)|[cut(`ip`)&cut(`port`)]\")\nFork(\"cut(`port`)\")\n",
		NewParser().MustParse("[ cut(`ip`) | [cut(`ip`) & cut(`port`)] & cut(`port`) ]"))
	assert.Equal(t, "Fork(\"cut(`ip`)|[cut(`ip`)&cut(\\\"port\\\")]\")\nFork(\"cut(`port`)\")\n",
		NewParser().MustParse("[ cut(`ip`) | [cut(`ip`) & cut(\"port\")] & cut(`port`) ]"))

	// fork的语法错误
	_, err := NewParser().Parse("[ cut(`ip`) & 1 & cut(`port`) ]")
	assert.Error(t, err)
}

func TestRegisterFunction(t *testing.T) {
	// translate simple mode to go code
	RegisterFunction("fofa", func(fi *FuncInfo) string {
		tmpl, err := template.New("fofa").Parse(`FetchFofa(GetRunner(), map[string]interface{} {
    "query": {{ .Query }},
    "size": {{ .Size }},
    "fields": {{ .Fields }},
})`)
		if err != nil {
			panic(err)
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, struct {
			Query  string
			Size   int64
			Fields string
		}{
			Query:  fi.Params[0].String(),
			Size:   fi.Params[1].Int64(),
			Fields: fi.Params[2].String(),
		})
		if err != nil {
			panic(err)
		}
		return tpl.String()
	})
	assert.Equal(t, `FetchFofa(GetRunner(), map[string]interface{} {
    "query": `+"`"+`title="test"`+"`"+`,
    "size": 10,
    "fields": `+"`"+`host,title,body`+"`"+`,
})
`, NewParser().MustParse("fofa(`title=\"test\"`, 10, `host,title,body`)"))
}

func TestParser_ParseToGraph(t *testing.T) {

	v, err := NewParser().ParseToGraph("test() | [ cut(`ip`) | [cut(`ip`) & cut(`port`)] & cut(`port`) ]")
	assert.Nil(t, err)
	assert.Equal(t, "graph TD\ntest1[\"test()\"]-->cut2[\"cut(`ip`)\"]\ncut2[\"cut(`ip`)\"]-->cut3[\"cut(`ip`)\"]\ncut2[\"cut(`ip`)\"]-->cut4[\"cut(`port`)\"]\ntest1[\"test()\"]-->cut5[\"cut(`port`)\"]\n", v)

	v, err = NewParser().ParseToGraph("fofa(`title=test`) | to_int(`port`) | sort(`port`) | [cut(`port`) & cut(`ip`)]")
	assert.Nil(t, err)
	assert.Equal(t, "graph TD\nfofa1[\"fofa(`title=test`)\"]-->to_int2[\"to_int(`port`)\"]\nto_int2[\"to_int(`port`)\"]-->sort3[\"sort(`port`)\"]\nsort3[\"sort(`port`)\"]-->cut4[\"cut(`port`)\"]\nsort3[\"sort(`port`)\"]-->cut5[\"cut(`ip`)\"]\n", v)

	v, err = NewParser().ParseToGraph(`cut("ip")|cut("port")`)
	assert.Nil(t, err)
	assert.Equal(t, "graph TD\ncut1[\"cut(#quot;ip#quot;)\"]-->cut2[\"cut(#quot;port#quot;)\"]\n", v)

	v, err = NewParser().ParseToGraph(`cut("ip")|cut("port")`, `graph LR`+"\n")
	assert.Nil(t, err)
	assert.Equal(t, "graph LR\ncut1[\"cut(#quot;ip#quot;)\"]-->cut2[\"cut(#quot;port#quot;)\"]\n", v)

}
