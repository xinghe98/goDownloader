package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// 获取文件总大小（通过 HTTP HEAD 请求）
func GetFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("服务器返回错误")
	}
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}
	fmt.Printf("文件总大小为: %.2fMB \n", float64(size/1024/1024))
	return size, nil
}
