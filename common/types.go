package common

import (
	"net/http"

	"github.com/vbauerster/mpb"
)

// worker接口
type Worker interface {
	start()
	stop()
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

// 工作池
type Pool struct {
	ThreadCount int
	Worker      *Worker
	Tasks
}
