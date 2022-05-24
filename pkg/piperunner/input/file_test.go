package input

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad_load(t *testing.T) {
	assert.Equal(t,
		`LoadFile(GetRunner(), map[string]interface{} {
    "file": "test.json",
})
`,
		pipeast.NewParser().Parse(`load("test.json")`))

	ast := pipeast.NewParser().Parse(`load("../../../data/forktest.json") | [cut("a") & cut("b")]`)
	p := piperunner.New()
	_, err := p.Run(ast)
	assert.Nil(t, err)
	res, err := os.ReadFile(p.LastFile)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,"b":2}`, string(res))
	assert.Equal(t, 2, len(p.LastTask.Children))
	res, err = os.ReadFile(p.LastTask.Children[0].GetLastFile())
	assert.Nil(t, err)
	assert.Equal(t, "{\"a\":1}\n", string(res))
	res, err = os.ReadFile(p.LastTask.Children[1].GetLastFile())
	assert.Nil(t, err)
	assert.Equal(t, "{\"b\":2}\n", string(res))
}
