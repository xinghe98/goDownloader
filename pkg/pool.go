package pkg

import (
	"runtime"

	"github.com/xinghe98/goDownloader/common"
)

// 初始化工作池
func NewPool(worker common.Worker, tasks common.Tasks) common.Pool {
	threads := runtime.NumCPU()
	return common.Pool{
		ThreadCount: threads,
		Worker:      &worker,
		Tasks:       tasks,
	}
}
