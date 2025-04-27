package downloadmp4

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/pkg"
)

func NewMp4Tasks(url string, outputFileName string, parts []*os.File, p *mpb.Progress) *mP4Tasks {
	threadCount := runtime.NumCPU()
	return &mP4Tasks{
		threadsCount: threadCount,
		url:          url,
		outputFile:   outputFileName,
		parts:        parts,
		p:            p,
	}
}

type mP4Tasks struct {
	threadsCount int
	url          string
	outputFile   string
	parts        []*os.File
	p            *mpb.Progress
}

func (mp4Task *mP4Tasks) EnterQueue(taskChan chan<- common.Tasks) {
	// 获取文件大小
	size, err := pkg.GetFileSize(mp4Task.url)
	if err != nil {
		println(err)
	}
	// 计算每个部分的大小
	partSize := size / int64(mp4Task.threadsCount)
	remainder := size % int64(mp4Task.threadsCount)

	// 将需要下载部分的req传入taskchan，后续使用
	for i := range mp4Task.threadsCount {
		// 创建临时文件部分
		partFile, err := os.Create(fmt.Sprintf("%s.part%d", mp4Task.outputFile, i))
		if err != nil {
			println(err)
		}
		mp4Task.parts[i] = partFile

		// 构造需要下载的请求
		start := int64(i) * partSize
		end := start + partSize - 1
		if i == mp4Task.threadsCount-1 {
			end += remainder // 最后一部分要加上多余的
		}
		req, err := http.NewRequest("GET", mp4Task.url, nil)
		if err != nil {
			println(err)
		}
		// 设置Range头
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

		// 添加一个进度条
		bar := mp4Task.p.AddBar(end-start,
			mpb.BarWidth(20), // 固定宽度为 20 字符
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("线程%s: ", strconv.Itoa(i+1))), // 进度条名称
				decor.Percentage(decor.WCSyncSpace),                  // 显示百分比
				decor.Name(" | "),                                    // 自定义分隔符
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
}
