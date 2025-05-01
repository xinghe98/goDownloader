package downloader

import (
	"strings"

	"github.com/xinghe98/goDownloader/common"
)

func NewDownloader(url string, outputName string) common.DownLoader {
	if strings.HasSuffix(url, "m3u8") {
		panic("no function")
	}
	return &Mp4Downloader{
		Url:        url,
		OutputFlie: outputName,
	}
}
