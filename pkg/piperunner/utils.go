package piperunner

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	defaultPipeTmpFilePrefix = "gofofa_pipeline_"
)

func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

// EachLine 每行处理文件
func EachLine(filename string, f func(line string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		f(string(line))
	}
	return nil
}

// WriteTempFile 写入临时文件
// 如果writeF是nil，就只返回生成的一个临时空文件路径
func WriteTempFile(ext string, writeF func(f *os.File)) string {
	var f *os.File
	var err error
	if len(ext) > 0 {
		ext = "*" + ext
	}
	f, err = os.CreateTemp(os.TempDir(), defaultPipeTmpFilePrefix+ext)
	if err != nil {
		panic(fmt.Errorf("create tmpfile failed: %w", err))
	}
	defer f.Close()

	if writeF != nil {
		writeF(f)
	}
	return f.Name()
}
