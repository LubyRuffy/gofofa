package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_scan_port(t *testing.T) {
	assert.Equal(t,
		`ScanPort(GetRunner(), map[string]interface{}{
    "targets": "127.0.0.1",
    "ports": "22,80,443,1080,3389,8000,8080,8443",
})
`,
		workflowast.NewParser().MustParse(`scan_port()`))

	assert.Equal(t,
		`ScanPort(GetRunner(), map[string]interface{}{
    "targets": "1.1.1.1",
    "ports": "22,80,443,1080,3389,8000,8080,8443",
})
`,
		workflowast.NewParser().MustParse(`scan_port("1.1.1.1")`))

	assert.Equal(t,
		`ScanPort(GetRunner(), map[string]interface{}{
    "targets": "1.1.1.1",
    "ports": "80",
})
`,
		workflowast.NewParser().MustParse(`scan_port("1.1.1.1","80")`))

}
