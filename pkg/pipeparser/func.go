package pipeparser

import (
	"strconv"
	"sync"
)

// FunctionTranslateHook 转换hook函数
type FunctionTranslateHook func(fi *FuncInfo) string

var (
	globalFunctionTranslateHooks sync.Map
)

type FuncParameter struct {
	v interface{}
}

// Int64 做为int64返回
func (fp FuncParameter) Int64() int64 {
	return fp.v.(int64)
}

// String 做为string返回
func (fp FuncParameter) String() string {
	return fp.v.(string)
}

// ToString 转换成字符串
func (fp FuncParameter) ToString() string {
	switch fp.v.(type) {
	case string:
		return fp.v.(string)
	case int64:
		return strconv.FormatInt(fp.v.(int64), 10)
	case *FuncInfo:
		return fp.v.(*FuncInfo).String()
	default:
		panic(fp.v)
	}
	return ""
}

// FuncInfo 函数信息
type FuncInfo struct {
	Name   string           // 函数名称
	Params []*FuncParameter // 参数列表
}

// String func id string
func (f *FuncInfo) String() string {

	if v, ok := globalFunctionTranslateHooks.Load(f.Name); ok {
		return v.(FunctionTranslateHook)(f)
	}

	rStr := f.Name + "("
	for i, p := range f.Params {
		if i != 0 {
			rStr += ", "
		}
		rStr += p.ToString()
	}
	return rStr + ")"
}

func RegisterFunction(name string, f FunctionTranslateHook) {
	globalFunctionTranslateHooks.Store(name, f)
}