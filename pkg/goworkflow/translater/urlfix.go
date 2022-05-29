package translater

import "github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"

func urlfixHook(fi *workflowast.FuncInfo) string {
	urlField := "url"
	if len(fi.Params) > 0 {
		urlField = fi.Params[0].RawString()
	}
	return `URLFix(GetRunner(), map[string]interface{}{
    "urlField": "` + urlField + `",
})`
}

func init() {
	register("urlfix", urlfixHook) // 补充完善url
}
