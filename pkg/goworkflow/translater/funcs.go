// Package translater workflow to gocode
package translater

import "github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"

var (
	Translators []string
)

// Load do nothing, only import translater
func register(name string, f workflowast.FunctionTranslateHook) {
	Translators = append(Translators, name)
	workflowast.RegisterFunction(name, f)
}
