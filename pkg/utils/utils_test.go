package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEachLine(t *testing.T) {
	var err error

	// 成功
	err = EachLine("../../data/forktest.json", func(line string) error {
		t.Logf(line)
		return nil
	})
	assert.Nil(t, nil)

	// 文件不存在
	err = EachLine("never_could_exists", nil)
	assert.True(t, os.IsNotExist(err))

	// 异常
	err = EachLine("../../data/forktest.json", func(line string) error {
		return errors.New("panic")
	})
	assert.Contains(t, err.Error(), "panic")
}

func TestWriteTempFile(t *testing.T) {
	writeContent := "abc"
	var f string
	var v []byte
	var err error

	writeF := func(f *os.File) error {
		f.WriteString(writeContent)
		return nil
	}

	// 正常，没有后缀
	f, err = WriteTempFile("", writeF)
	assert.Nil(t, err)
	assert.Contains(t, f, defaultPipeTmpFilePrefix)
	v, err = os.ReadFile(f)
	assert.Nil(t, err)
	assert.Contains(t, f, defaultPipeTmpFilePrefix)
	assert.Equal(t, writeContent, string(v))

	// 正常，有后缀
	f, err = WriteTempFile(".txt", writeF)
	assert.Equal(t, ".txt", filepath.Ext(f))

	// 异常
	f, err = WriteTempFile("/../../../../../../.txt", writeF)
	assert.Error(t, err) // os.errPatternHasSeparator
}

func TestEscapeString(t *testing.T) {
	assert.Equal(t, `\"\"`, EscapeString(`""`))
	assert.Equal(t, `#quot;#quot;`, EscapeDoubleQuoteStringOfHTML(`""`))
}

func TestReadFirstLineOfFile(t *testing.T) {
	f, err := WriteTempFile("", func(f *os.File) error {
		_, err := f.WriteString("aaaa\nbbbb\ncccc")
		return err
	})
	assert.Nil(t, err)

	data, err := ReadFirstLineOfFile("nevercouldexists")
	assert.Error(t, err)

	data, err = ReadFirstLineOfFile(f)
	assert.Nil(t, err)
	assert.Equal(t, "aaaa", string(data))
}

func TestJSONLineFields(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, JSONLineFields(`{"a":1,"b":2}`))
}

func TestFileExists(t *testing.T) {
	assert.True(t, FileExists("../../data/forktest.json"))
	assert.False(t, FileExists("nevercouldexists"))
}

func TestDockerStatusOk(t *testing.T) {
	defer func() {
		// 不要影响后续的测试
		defaultDockerPath = "docker"
		lastCheckDockerTime = time.Now().Add(-10 * defaultCheckDockerDuration)
	}()

	defaultDockerPath = "nevercouldexists"
	assert.False(t, DockerStatusOk())
	assert.False(t, DockerStatusOk())

	lastCheckDockerTime = time.Now().Add(-10 * defaultCheckDockerDuration)
	defaultDockerPath = "curl"
	assert.False(t, DockerStatusOk())
	assert.False(t, DockerStatusOk())

	// 跳过缓存
	lastCheckDockerTime = time.Now().Add(-10 * defaultCheckDockerDuration)
	defaultDockerPath = "docker"
	if DockerStatusOk() {
		// 看似没有意义，但是不确定测试环境有docker
		assert.True(t, DockerStatusOk())
	}

	// 跳过缓存正常
	lastCheckDockerTime = time.Now().Add(-10 * defaultCheckDockerDuration)
	defaultDockerPath = "docker"
	if DockerStatusOk() {
		// 看似没有意义，但是不确定测试环境有docker
		assert.True(t, DockerStatusOk())

		d, err := DockerRun("images")
		assert.Nil(t, err)
		assert.Contains(t, string(d), "REPOSITORY")
	}
}

func TestExecCmdWithTimeout(t *testing.T) {
	d, err := ExecCmdWithTimeout(time.Microsecond, "whoami")
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Nil(t, d)
	d, err = ExecCmdWithTimeout(3*time.Second, "whoami")
	assert.Nil(t, err)
	assert.True(t, len(d) > 0)
}

func TestSimpleHash(t *testing.T) {
	assert.Equal(t, "0x340ca71c", SimpleHash("1"))
}
