package utils

import (
	"bufio"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
	"hash/fnv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	defaultPipeTmpFilePrefix   = "gofofa_pipeline_"
	lastCheckDockerTime        time.Time // 最后检查docker路径的时间
	defaultDockerPath          = "docker"
	defaultCheckDockerDuration = 5 * time.Minute
	globalDockerOK             = false
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

// JSONLineFieldsWithType 获取json行的fields，包含属性信息
func JSONLineFieldsWithType(line string) (fields [][]string) {
	v := gjson.Parse(line)
	v.ForEach(func(key, value gjson.Result) bool {
		typeStr := "text"
		switch value.Type {
		case gjson.True, gjson.False:
			typeStr = "boolean"
		case gjson.Number:
			typeStr = "int"
		}
		fields = append(fields, []string{key.String(), typeStr})
		return true
	})
	return
}

// JSONLineFields 获取json行的fields
func JSONLineFields(line string) (fields []string) {
	fs := JSONLineFieldsWithType(line)
	for _, f := range fs {
		fields = append(fields, f[0])
	}
	return
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// LoadFirstExistsFile 从文件列表中返回第一个存在的文件路径
func LoadFirstExistsFile(paths []string) string {
	for _, p := range paths {
		if FileExists(p) {
			return p
		}
	}
	return ""
}

//// GetCurrentProcessFileDir 获得当前程序所在的目录
//func GetCurrentProcessFileDir() string {
//	return filepath.Dir(os.Args[0])
//}
//
//// UserHomeDir 获得当前用户的主目录
//func UserHomeDir() string {
//	if runtime.GOOS == "windows" {
//		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
//		if home == "" {
//			home = os.Getenv("USERPROFILE")
//		}
//		return home
//	}
//	return os.Getenv("HOME")
//}

// ExecCmdWithTimeout 在时间范围内执行系统命令，并且将输出返回（stdout和stderr）
func ExecCmdWithTimeout(timeout time.Duration, arg ...string) (b []byte, err error) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	routeCmd := exec.CommandContext(ctx, arg[0], arg[1:]...)

	return routeCmd.CombinedOutput()
}

// RunCmdNoExitError 将exec.ExitError不作为错误，通常配合exec.Command使用
func RunCmdNoExitError(d []byte, err error) ([]byte, error) {
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			err = nil
		}
	}
	return d, err
}

// DockerRun 运行docker，解决Windows找不到的问题
// 注意：exec.ExitError 错误会被忽略，我们只关心所有的字符串返回，不关注进程的错误代码
func DockerRun(args ...string) ([]byte, error) {
	// 缓存5分钟
	if time.Now().Sub(lastCheckDockerTime) < defaultCheckDockerDuration {
		if globalDockerOK {
			return RunCmdNoExitError(exec.Command(defaultDockerPath, args...).CombinedOutput())
		} else {
			return nil, fmt.Errorf("docker status is not ok")
		}
	}
	lastCheckDockerTime = time.Now()

	d, err := RunCmdNoExitError(exec.Command(defaultDockerPath, "version").CombinedOutput())
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			err = nil
		} else {
			// 可能路径不在PATH环境变量，需要自己找，主要是windows
			// https://docs.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables
			defaultDockerPath = LoadFirstExistsFile([]string{
				"docker.exe",
				filepath.Join(os.Getenv("PROGRAMFILES"), "Docker", "Docker", "resources", "bin", "docker.exe"),
				filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Docker", "Docker", "resources", "bin", "docker.exe"),
			})
			if len(defaultDockerPath) == 0 {
				return nil, fmt.Errorf("could not find docker")
			}
			d, err = RunCmdNoExitError(exec.Command(defaultDockerPath, "version").CombinedOutput())
		}
	}
	if err == nil {
		if strings.Contains(string(d), "API version") {
			globalDockerOK = true
			return RunCmdNoExitError(exec.Command(defaultDockerPath, args...).CombinedOutput())
		} else {
			err = fmt.Errorf("docker is invalid")
		}
	}

	return nil, err
}

// DockerStatusOk 检查是否安装
func DockerStatusOk() bool {
	data, err := DockerRun("images")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "REPOSITORY")
}

// SimpleHash hashes using fnv32a algorithm
func SimpleHash(text string) string {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return fmt.Sprintf("0x%08x", algorithm.Sum32())
}
