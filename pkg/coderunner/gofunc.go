package coderunner

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

// GoFunction user defined functions
type GoFunction struct {
	functions sync.Map
}

// Range 遍历用户自定义的函数
func (u *GoFunction) Range(f func(key, value any) bool) {
	u.functions.Range(f)
}

// Register 注册底层函数
func (u *GoFunction) Register(key string, f interface{}) error {
	// 格式
	if !regexp.MustCompile(`^[A-Za-z][0-9a-zA-Z_]*$`).MatchString(key) {
		return fmt.Errorf("function name is invalid: %s", key)
	}

	fType := reflect.ValueOf(f).Kind().String()
	if fType != "func" {
		return fmt.Errorf("function is invalid: %v", f)
	}

	u.functions.Store(key, f)
	return nil
}
