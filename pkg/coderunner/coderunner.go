// Package coderunner 底层代码的执行器（go语言）
// 最纯净的版本，不做任何多余的动作
package coderunner

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
)

// Runner simplest script runner
type Runner struct {
	functions *GoFunction // user defined functions
}

// Run code pipelines
func (p *Runner) Run(code string) (reflect.Value, error) {
	var err error
	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)

	exports := interp.Exports{
		"this/this": {
			"ENGINE":    reflect.ValueOf("runner"),
			"GetRunner": reflect.ValueOf(p), // 上层可以替换
		},
	}

	if p.functions != nil {
		p.functions.Range(func(key, value any) bool {
			exports["this/this"][key.(string)] = reflect.ValueOf(value)
			return true
		})
	}

	err = i.Use(exports)
	if err != nil {
		panic(err)
	}

	i.ImportUsed()
	i.Eval(`import (
	. "this/this"
)`)

	return i.Eval(code)
}

// New create go runner
func New(options ...RunnerOption) *Runner {
	r := &Runner{}
	for _, o := range options {
		o(r)
	}
	return r
}
