package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
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
