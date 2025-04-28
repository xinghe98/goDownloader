package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/vbauerster/mpb"
	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/downloadmp4"
	"github.com/xinghe98/goDownloader/pkg"
)

type Options struct {
	Url    string `short:"u" long:"url" description:"下载链接的URL" required:"true"`
	Output string `short:"o" long:"name" description:"导出的文件名" default:"output.mp4"`
}

func run(url string, outputFile string) {
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

func main() {
	// url := "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_30mb.mp4"
	// outputFile := "video.mp4"
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		// 处理错误
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrRequired {
			os.Exit(1)
			return
		}
		return
	}
	run(opts.Url, opts.Output)
}
