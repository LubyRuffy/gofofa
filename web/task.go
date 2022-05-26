package web

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	globalTaskMonitor = newTaskMonitor()
)

type taskInfo struct {
	monitor *taskMonitor
	taskId  string
	msgCh   chan string
	started time.Time
	ended   time.Time
	html    string
}

func (t *taskInfo) finish() {
	t.ended = time.Time{}
	close(t.msgCh)
	go func() {
		select {
		case <-time.After(time.Minute):
			// 1分钟后删除
			t.monitor.del(t.taskId)
		}
	}()
}

func (t *taskInfo) addMsg(msg string) {
	t.msgCh <- msg
}

func (t *taskInfo) receiveMsg() (string, bool) {
	select {
	case msg, ok := <-t.msgCh:
		return msg, ok
	}
	return "", false
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

func (t *taskMonitor) new() *taskInfo {
	tid := uuid.New().String()
	ti := &taskInfo{
		taskId:  tid,
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
