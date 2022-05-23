package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func assertPipeCmd(t *testing.T, pipeCmd string, jsonData string, except string) {
	p := piperunner.New(pipeparser.NewParser().Parse(pipeCmd))

	// write json to file
	f, err := os.CreateTemp(os.TempDir(), "piperunner_")
	assert.Nil(t, err)
	defer f.Close()
	_, err = f.WriteString(jsonData)
	assert.Nil(t, err)

	p.LastFile = f.Name()

	// run
	err = p.Run()
	assert.Nil(t, err)

	res, err := os.ReadFile(p.LastFile)
	assert.Nil(t, err)

	assert.Equal(t, except, string(res))
}

func TestPipeRunner_flat(t *testing.T) {
	assertPipeCmd(t, `flat("a")`, `{"a":[[1], "a", [[[true]]]}`, "1\n\"a\"\ntrue\n")
}
