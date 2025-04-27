package downloadmp4

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/pkg"
)

func NewMp4Worker(part []*os.File, taskChan <-chan common.Tasks, resultChan chan<- common.Resluts, wg *sync.WaitGroup) *downloadMP4Worker {
	return &downloadMP4Worker{
		parts:      part,
		taskChan:   taskChan,
		resultChan: resultChan,
		wg:         wg,
	}
}

type downloadMP4Worker struct {
	parts      []*os.File
	taskChan   <-chan common.Tasks
	resultChan chan<- common.Resluts
	wg         *sync.WaitGroup
}

func (downloadMP4 *downloadMP4Worker) Start() {
	client := pkg.GetClient()
	defer downloadMP4.wg.Done()
	for task := range downloadMP4.taskChan {
		result := common.Resluts{}
		resp, err := client.Do(task.Req)
		if err != nil {
			result.Err = err
			downloadMP4.resultChan <- result
			continue
		}

		if resp.StatusCode != http.StatusPartialContent {
			result.Err = fmt.Errorf("服务器不支持分段下载，返回状态: %s", resp.Status)
			downloadMP4.resultChan <- result
			resp.Body.Close()
			continue
		}
		func() {
			// 使用 ProxyReader 自动更新进度条
			reader := task.Bar.ProxyReader(resp.Body)
			_, err = io.Copy(downloadMP4.parts[task.Index], reader)
			// fmt.Printf("开始下载第%s个片段,大小：%.2fMB\n", strconv.Itoa(task.Index), float64(task.Size/1024/1024))
			if err != nil {
				result.Err = err
			}
			defer resp.Body.Close()
			defer reader.Close()
		}()
		result.Err = nil
		result.Success = true
		downloadMP4.resultChan <- result
		task.Bar.Increment()
	}
}
