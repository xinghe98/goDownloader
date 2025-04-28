package common

import (
	"net/http"
	"sync"

	"github.com/vbauerster/mpb"
)

// worker接口
type Worker interface {
	Start(taskChan <-chan Tasks, resultChan chan<- Resluts, wg *sync.WaitGroup)
}

// Tasker接口
type Tasker interface {
	EnterQueue(taskChan chan<- Tasks)
}

// pool接口
type Pool interface {
	Start(taskChan chan<- Tasks, resultChan <-chan Resluts, wg *sync.WaitGroup)
}

// 任务队列
type Tasks struct {
	Index int
	Req   *http.Request
	Size  int64
	Bar   *mpb.Bar
}

// 结果队列
type Resluts struct {
	Err     error
	Success bool
}
