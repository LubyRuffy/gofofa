package utils

import (
	"bufio"
	"github.com/tidwall/gjson"
	"io"
	"os"
	"strings"
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
				return err
			}
		}

		err = f(string(line))
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteTempFile 写入临时文件
// 如果writeF是nil，就只返回生成的一个临时空文件路径
// 返回文件名和错误
func WriteTempFile(ext string, writeF func(f *os.File) error) (fn string, err error) {
	var f *os.File
	if len(ext) > 0 {
		ext = "*" + ext
	}
	f, err = os.CreateTemp(os.TempDir(), defaultPipeTmpFilePrefix+ext)
	if err != nil {
		return
	}
	defer f.Close()

	fn = f.Name()

	if writeF != nil {
		err = writeF(f)
		if err != nil {
			return
		}
	}
	return
}

// EscapeString 双引号内的字符串转换
func EscapeString(s string) string {
	//s, _ = sjson.Set(`{"a":""}`, "a", s)
	//return s[strings.Index(s, `:`)+1 : len(s)-1]
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// EscapeDoubleQuoteStringOfHTML 双引号内的字符串转换为Mermaid格式（HTML）
func EscapeDoubleQuoteStringOfHTML(s string) string {
	s = strings.ReplaceAll(s, `"`, `#quot;`)
	return s
}

// ReadFirstLineOfFile 读取文件的第一行
func ReadFirstLineOfFile(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var b [1]byte
	var data []byte
	for {
		_, err = f.Read(b[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, err
		}
		if b[0] == '\n' {
			break
		}
		data = append(data, b[0])
	}
	return data, nil
}

// JSONLineFields 获取json行的fields
func JSONLineFields(line string) (fields []string) {
	v := gjson.Parse(line)
	v.ForEach(func(key, value gjson.Result) bool {
		fields = append(fields, key.String())
		return true
	})
	return
}
