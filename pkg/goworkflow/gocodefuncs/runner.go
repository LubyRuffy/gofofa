package gocodefuncs

import (
	"github.com/lubyruffy/gofofa"
	"github.com/sirupsen/logrus"
)

type Runner interface {
	GetLastFile() string        // GetLastFile 获取最后一次生成的文件
	GetFofaCli() *gofofa.Client // GetFofaCli 获取fofa客户端连接

	Debugf(format string, args ...interface{})                   // 打印调试信息
	Warnf(format string, args ...interface{})                    // 打印警告信息
	Logf(level logrus.Level, format string, args ...interface{}) // 打印日志信息
}

// Artifact 过程中生成的文件
type Artifact struct {
	FilePath string // 文件路径
	FileName string // 文件路径
	FileSize int    // 文件大小
	FileType string // 文件类型
	Memo     string // 备注，比如URL等
}

// FuncResult 返回的结构
type FuncResult struct {
	OutFile   string // 往后传递的文件
	Artifacts []*Artifact
}
