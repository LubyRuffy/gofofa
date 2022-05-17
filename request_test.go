package gofofa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Fetch(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:55")
	assert.Error(t, err)
}
