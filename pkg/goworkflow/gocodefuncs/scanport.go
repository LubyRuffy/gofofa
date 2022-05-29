package gocodefuncs

import (
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"os"
	"strings"
)

type scanPortParam struct {
	Targets string `json:"targets"` // 扫描目标
	Ports   string `json:"ports"`   // 扫描端口
}

// ScanPort 扫描端口
// 参数: hosts/ports
// 输出格式：{"ip":"117.161.125.154","port":80,"protocol":"tcp","service":"http","hostnames":"fofa.info"}
func ScanPort(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options scanPortParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("screenShot failed: %w", err))
	}

	if options.Targets == "" {
		options.Targets = "127.0.0.1"
	}
	if options.Ports == "" {
		options.Ports = "22,80,443,1080,3389,8080,8443"
	}

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(options.Targets),
		nmap.WithPorts(options.Ports),
		nmap.WithOpenOnly(),
	)
	if err != nil {
		panic(fmt.Errorf("ScanPort error: %w", err))
	}

	progress := make(chan float32, 1)
	// Function to listen and print the progress
	go func() {
		for p := range progress {
			fmt.Printf("Progress: %v %%\n", p)
		}
	}()

	result, _, err := scanner.RunWithProgress(progress)
	if err != nil {
		panic(fmt.Errorf("ScanPort error: %w", err))
	}

	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		// 遍历host
		for _, host := range result.Hosts {
			if len(host.Ports) == 0 || len(host.Addresses) == 0 {
				continue
			}
			var hostnames []string
			for i := range host.Hostnames {
				hostnames = append(hostnames, host.Hostnames[i].Name)
			}
			for _, addr := range host.Addresses {
				for _, port := range host.Ports {
					_, err := f.WriteString(fmt.Sprintf(`{"ip":"%s","port":%d,"protocol":"%s","service":"%s","hostnames":"%s"}`+"\n",
						addr, port.ID, port.Protocol, port.Service.Name, strings.Join(hostnames, ",")))
					if err != nil {
						return err
					}
				}
			}
		}

		return err
	})
	if err != nil {
		panic(fmt.Errorf("ScanPort error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
