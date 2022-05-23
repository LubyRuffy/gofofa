package funcs

import (
	"testing"
)

func TestPipeRunner_cut(t *testing.T) {
	assertPipeCmd(t, `cut("a")`, `{"a":1,"b":2}`, "{\"a\":1}\n")
}
