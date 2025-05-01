package downloadmp4

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/xinghe98/goDownloader/common"
	"github.com/xinghe98/goDownloader/pkg"
)

func NewMp4Worker(part []*os.File) *downloadMP4Worker {
	return &downloadMP4Worker{
		parts: part,
	}
}

type downloadMP4Worker struct {
	parts []*os.File
}

func (downloadMP4 *downloadMP4Worker) Start(taskChan <-chan common.Tasks, resultChan chan<- common.Resluts) {
	client := pkg.GetClient()
	for task := range taskChan {
		result := common.Resluts{}
		resp, err := client.Do(task.Req)
		if err != nil {
			result.Err = err
			resultChan <- result
			continue
		}

		if resp.StatusCode != http.StatusPartialContent {
			result.Err = fmt.Errorf("服务器不支持分段下载，返回状态: %s", resp.Status)
			resultChan <- result
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
		resultChan <- result
		task.Bar.Increment()
	}
}
