package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/pkg"
)

func main() {
	start := time.Now()
	url := "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_30mb.mp4"
	// url := "https://1470423830.rsc.cdn77.org/mp43/1077727.mp4?secure=umfCgFZB0i17qlvjUIxCpg==,1745159423&f=ed94CZ5bhzdVrloADdPvsKhYK4RrulMCtNFZFQJvocxjgvlwwACwxjumS3/DxUQooEk+H43uSDFE7Q1BZ8SeGMyIbkLN2OazKibJjZY"
	outputFile := "video.mp4"

	// 定义工作者数量
	threads := runtime.NumCPU()
	parts := make([]*os.File, threads)
	var wg sync.WaitGroup

	// 创建进度条容器
	p := mpb.New()

	// 初始化通道
	taskChan := make(chan common.Tasks, threads)
	resultChan := make(chan common.Resluts, threads)

	// 启动工作者
	for range threads {
		wg.Add(1)
		go pkg.DownloadFile(parts, taskChan, resultChan, &wg)
	}

	// 获取文件大小
	size, err := pkg.GetFileSize(url)
	if err != nil {
		println(err)
	}
	// 计算每个部分的大小
	partSize := size / int64(threads)
	remainder := size % int64(threads)

	// 将需要下载部分的req传入taskchan，后续使用
	for i := range threads {
		// 创建临时文件部分
		partFile, err := os.Create(fmt.Sprintf("%s.part%d", outputFile, i))
		if err != nil {
			println(err)
		}
		parts[i] = partFile

		// 构造需要下载的请求
		start := int64(i) * partSize
		end := start + partSize - 1
		if i == threads-1 {
			end += remainder // 最后一部分要加上多余的
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			println(err)
		}
		// 设置Range头
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

		// 添加一个进度条，总任务量为 100
		bar := p.AddBar(end-start,
			mpb.BarWidth(20), // 固定宽度为 20 字符
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("线程%s: ", strconv.Itoa(i))), // 进度条名称
				decor.Percentage(decor.WCSyncSpace),                // 显示百分比
				decor.Name(" | "),                                  // 自定义分隔符
				decor.Counters(decor.UnitKiB, "%.1f/%.1f"),
			),
			mpb.AppendDecorators(
				decor.AverageSpeed(decor.UnitKiB, "%.1f"), // 显示 MB/s
			),
		)
		task := &common.Tasks{
			Index: i,
			Size:  end - start,
			Req:   req,
			Bar:   bar,
		}
		taskChan <- *task
	}
	close(taskChan)

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
		}
	}

	pkg.MergeParts(outputFile, parts)
}
