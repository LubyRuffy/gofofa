package utils

import (
	"bufio"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	defaultPipeTmpFilePrefix = "gofofa_pipeline_"
	lastCheckDockerTime      time.Time // 最后检查docker路径的时间
	dockerPath               = "docker"
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

// GetCurrentProcessFileDir 获得当前程序所在的目录
func GetCurrentProcessFileDir() string {
	return filepath.Dir(os.Args[0])
}

// UserHomeDir 获得当前用户的主目录
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ExecCmdWithTimeout 在时间范围内执行系统命令，并且将输出返回（stdout和stderr）
func ExecCmdWithTimeout(timeout time.Duration, arg ...string) (b []byte, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("[WARNING] ExecCmdWithTimeout failed: %v", x)
		}
	}()

	routeCmd := exec.Command(arg[0], arg[1:]...)

	var timer *time.Timer
	timer = time.AfterFunc(timeout, func() {
		timer.Stop()
		if routeCmd.Process != nil {
			routeCmd.Process.Kill()
		}
	})

	return routeCmd.CombinedOutput()
}

// DockerRun 运行docker，解决Windows找不到的问题
func DockerRun(args ...string) ([]byte, error) {
	// 缓存5分钟
	if time.Now().Sub(lastCheckDockerTime) > 5*time.Minute {
		return exec.Command(dockerPath, args...).CombinedOutput()
	}

	dockerPath = "docker"
	d, err := exec.Command(dockerPath, "version").CombinedOutput()
	if err != nil {
		// 可能路径不在PATH环境变量，需要自己找
		// https://docs.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables
		dockerPath := LoadFirstExistsFile([]string{
			"docker.exe",
			filepath.Join(os.Getenv("PROGRAMFILES"), "Docker", "Docker", "resources", "bin", "docker.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Docker", "Docker", "resources", "bin", "docker.exe"),
		})
		if len(dockerPath) > 0 {
			dockerPath = dockerPath
		}
		d, err = exec.Command(dockerPath, "version").CombinedOutput()
	}
	if err == nil {
		if strings.Contains(string(d), "API version") {
			return exec.Command(dockerPath, args...).CombinedOutput()
		}
	}

	return nil, err
}

// DockerStatusOk 检查是否安装
func DockerStatusOk() bool {
	_, err := DockerRun("version")
	return err == nil
}
