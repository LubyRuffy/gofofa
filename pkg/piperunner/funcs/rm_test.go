package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_rm(t *testing.T) {
	assert.Equal(t,
		`RemoveField(GetRunner(), map[string]interface{}{
   "fields": "title",
})
`,
		pipeast.NewParser().Parse(`rm("title")`))

	assertPipeCmd(t, `rm("title")`,
		`{"title":"abc","a":1}`,
		`{"a":1}
`)
}
