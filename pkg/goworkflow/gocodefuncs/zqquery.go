package gocodefuncs

import (
	"fmt"
	"github.com/brimdata/zed/cli/zq"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type zqQueryParams struct {
	Query string `json:"query"`
}

// ZqQuery zq command
func ZqQuery(p Runner, params map[string]interface{}) *FuncResult {
	var fn string
	var err error
	var options zqQueryParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	fn, err = utils.WriteTempFile(".json", nil)
	if err != nil {
		panic(err)
	}

	//err = fzq.ZqQuery(options.Query, p.GetLastFile(), fn)
	cmd := []string{"-f", "json", "-o", fn, options.Query, p.GetLastFile()}
	logrus.Debugf("zq cmd: %v", cmd)
	err = zq.Cmd.ExecRoot(cmd)
	if err != nil {
		panic(fmt.Errorf("zqQuery error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
