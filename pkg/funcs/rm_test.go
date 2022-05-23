package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_rm(t *testing.T) {
	assert.Equal(t,
		`RemoveField(GetRunner(), map[string]interface{}{
   "fields": "title",
})
`,
		pipeparser.NewParser().Parse(`rm("title")`))

	assertPipeCmd(t, `rm("title")`,
		`{"title":"abc","a":1}`,
		`{"a":1}
`)
}
