package gocodefuncs

import (
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"os"
	"path/filepath"
)

type loadFileParams struct {
	File string
}

// LoadFile 加载json文件
func LoadFile(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options loadFileParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("loadFile failed: %w", err))
	}

	if len(options.File) == 0 {
		panic(errors.New("load file cannot be empty"))
	}

	var path string
	//path, _ = os.Getwd()
	path, _ = filepath.Abs(options.File)

	if _, err = os.Stat(path); err != nil {
		panic(fmt.Errorf("loadFile failed: %w", err))
	}

	//return path, nil

	fn, err := utils.WriteTempFile(".json", func(f *os.File) error {
		var bytesRead []byte
		bytesRead, err = ioutil.ReadFile(options.File)
		if err != nil {
			panic(err)
		}
		_, err = f.Write(bytesRead)
		return err
	})

	if err != nil {
		panic(fmt.Errorf("loadFile error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
