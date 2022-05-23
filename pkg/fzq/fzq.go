// Package fzq fofa zq
package fzq

import (
	"github.com/brimdata/zed/cli/zq"
)

func ZqQuery(zed string, inputFile string, outputFile string) error {
	cmd := []string{"-f", "json", "-o", outputFile, zed, inputFile}
	return zq.Cmd.ExecRoot(cmd)
}
