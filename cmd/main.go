package main

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

func main() {
	// INFO: 先全局定义一些内容
	url := "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_30mb.mp4"
	start := time.Now()
	outputFile := "video.mp4"
	threads := runtime.NumCPU()
	parts := make([]*os.File, threads)
	var wg sync.WaitGroup
	// 创建进度条容器
	p := mpb.New()

	// 初始化通道
	taskChan := make(chan common.Tasks, threads)
	resultChan := make(chan common.Resluts, threads)

	// 初始化工作者
	worker := downloadmp4.NewMp4Worker(parts, taskChan, resultChan, &wg)
	// 初始化任务队列
	tasker := downloadmp4.NewMp4Tasks(url, outputFile, parts, p)
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

	pkg.MergeParts(outputFile, parts)
}
