package translater

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	assert.Equal(t, 17, len(Translators))
}
