package downloader

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/downloadmp4"
	"github.com/xinghe98/goDownloader/pkg"
)

type Mp4Downloader struct {
	Url        string
	OutputFlie string
}

func (Mp4 *Mp4Downloader) Download() {
	start := time.Now()
	threads := runtime.NumCPU()
	parts := make([]*os.File, threads)
	var wg sync.WaitGroup
	// 创建进度条容器
	p := mpb.New()

	// 初始化通道
	taskChan := make(chan common.Tasks, threads)
	resultChan := make(chan common.Resluts, threads)

	// 初始化工作者
	worker := downloadmp4.NewMp4Worker(parts)
	// 初始化任务队列
	tasker := downloadmp4.NewMp4Tasks(Mp4.Url, Mp4.OutputFlie, parts, p)
	// 初始化工作池
	pool := pkg.NewPool(worker, tasker)
	// 启动工作池
	pool.Start(taskChan, resultChan, &wg)

	go func() {
		wg.Wait()
		p.Wait()
		close(resultChan)
		cost := time.Since(start)
		fmt.Printf("总耗时：[%s]\n", cost)
	}()

	// 处理结果，合并为mp4
	for result := range resultChan {
		if result.Err != nil {
			fmt.Printf("❌ Task failed: %s \n", result.Err.Error())
			return
		}
	}
	pkg.MergeParts(Mp4.OutputFlie, parts)
}
