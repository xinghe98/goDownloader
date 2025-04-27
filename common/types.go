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

type Tasker interface {
	EnterQueue(taskChan chan<- Tasks)
}

// pool接口
type Pool interface {
	Start(taskChan chan<- Tasks, resultChan <-chan Resluts, wg *sync.WaitGroup)
}

type Tasks struct {
	Index int
	Req   *http.Request
	Size  int64
	Bar   *mpb.Bar
}

type Resluts struct {
	Err     error
	Success bool
}
