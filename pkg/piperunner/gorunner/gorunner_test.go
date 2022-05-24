package gorunner

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	gr := New()
	v, err := gr.Run("ENGINE")
	assert.Nil(t, err)
	assert.Equal(t, "runner", v.String())

	_, err = gr.Run("MyFunc()")
	assert.Error(t, err)
}

func TestWithFunctions(t *testing.T) {
	var err error
	var v reflect.Value

	gf := GoFunction{}
	err = gf.Register("MyFunc", func() string {
		return "hello"
	})
	gr := New(WithFunctions(&gf))
	v, err = gr.Run("MyFunc()")
	assert.Nil(t, err)
	assert.Equal(t, "hello", v.String())
}
