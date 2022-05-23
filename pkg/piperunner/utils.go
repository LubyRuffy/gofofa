package piperunner

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

// WriteTempJSONFile 写入临时文件
func WriteTempJSONFile(writeF func(f *os.File)) string {
	var f *os.File
	var err error
	f, err = os.CreateTemp(os.TempDir(), defaultPipeTmpFilePrefix)
	if err != nil {
		panic(fmt.Errorf("create tmpfile failed: %w", err))
	}
	defer f.Close()

	writeF(f)

	return f.Name()
}
