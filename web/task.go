package web

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	globalTaskMonitor = newTaskMonitor()
)

type taskInfo struct {
	monitor       *taskMonitor
	runner        *goworkflow.PipeRunner
	taskId        string
	code          string // 运行的代码
	msgCh         chan string
	started       time.Time
	ended         time.Time
	html          string
	callIDRunning int // 当前运行的callID
	finished      bool
}

func (t *taskInfo) finish() {
	t.ended = time.Time{}
	t.finished = true
	close(t.msgCh)

	go func() {
		select {
		case <-time.After(1 * time.Minute):
			// 1分钟后删除
			t.monitor.del(t.taskId)
		}
	}()
}

func (t *taskInfo) addMsg(msg string) {
	if t.finished {
		return
	}
	t.msgCh <- msg
}

func (t *taskInfo) receiveMsg() (string, bool) {
	select {
	case msg, ok := <-t.msgCh:
		return msg, ok
	default:
		return "", false
	}
}

type taskMonitor struct {
	m sync.Map
}

func newTaskMonitor() *taskMonitor {
	return &taskMonitor{}
}

func (t *taskMonitor) del(taskId string) {
	t.m.Delete(taskId)
}

func (t *taskMonitor) new(code string) *taskInfo {
	tid := uuid.New().String()
	ti := &taskInfo{
		taskId:  tid,
		code:    code,
		msgCh:   make(chan string, 1000),
		started: time.Now(),
		monitor: t,
	}
	t.m.Store(tid, ti)
	return ti
}

func (t *taskMonitor) addMsg(taskid string, msg string) {
	if task, ok := t.m.Load(taskid); ok {
		task.(*taskInfo).addMsg(msg)
	}
}

func (t *taskMonitor) receiveMsg(taskid string) (string, bool) {
	if task, ok := t.m.Load(taskid); ok {
		return task.(*taskInfo).receiveMsg()
	}
	return "", false
}
