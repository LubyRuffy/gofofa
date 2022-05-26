package translater

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	assert.Equal(t, 19, len(Translators))
}
