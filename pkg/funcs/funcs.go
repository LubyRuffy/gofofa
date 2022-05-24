package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeast"

// SupportWorkflows 手动加载，否则init不执行
func SupportWorkflows() []string {
	return pipeast.SupportWorkflows()
}
