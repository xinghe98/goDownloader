package pkg

import (
	"runtime"
	"sync"

	"github.com/xinghe98/goDownloader/common"
)

// 初始化工作池
func NewPool(worker common.Worker, tasks common.Tasker) pool {
	threads := runtime.NumCPU()
	return pool{
		ThreadCount: threads,
		Worker:      worker,
		Tasks:       tasks,
	}
}

// 工作池
type pool struct {
	ThreadCount int
	Worker      common.Worker
	Tasks       common.Tasker
}

func (p *pool) Start(taskChan chan common.Tasks, resultChan chan common.Resluts, wg *sync.WaitGroup) {
	// 启动工作者
	for range p.ThreadCount {
		wg.Add(1)
		go p.Worker.Start(taskChan, resultChan, wg)
	}

	// 构造任务，放入任务队列（通道）提供给工作者们
	p.Tasks.EnterQueue(taskChan)
}
