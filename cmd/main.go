package main

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/xinghe98/goDownloader/cmd/downloader"
	"github.com/xinghe98/goDownloader/common"
)

func NewDownloader(url string, outputName string) common.DownLoader {
	if strings.HasSuffix(url, "m3u8") {
		panic("no function")
	}
	return &downloader.Mp4Downloader{
		Url:        url,
		OutputFlie: outputName,
	}
}

type Options struct {
	Url    string `short:"u" long:"url" description:"下载链接的URL" required:"true"`
	Output string `short:"n" long:"name" description:"导出的文件名" default:"output.mp4"`
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
	downloader := NewDownloader(opts.Url, opts.Output)
	downloader.Download()
}
