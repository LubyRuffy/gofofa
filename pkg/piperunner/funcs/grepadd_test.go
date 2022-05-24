package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_grepAdd(t *testing.T) {
	assert.Equal(t,
		"AddField(GetRunner(), map[string]interface{}{\n    \"from\": map[string]interface{}{\n        \"method\": \"grep\",\n        \"field\": \"title\",\n        \"value\": \"(?is)test\",\n    },\n    \"name\": \"new_title\",\n})\n",
		pipeast.NewParser().Parse(`grep_add("title", "(?is)test", "new_title")`))

	assertPipeCmd(t, `grep_add("title", "(?is)test", "new_title")`,
		`{"title":"Test123"}
{"title":"123test456"}`,
		`{"title":"Test123","new_title":[["Test"]]}
{"title":"123test456","new_title":[["test"]]}`)
}
