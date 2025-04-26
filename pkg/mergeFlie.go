package pkg

import (
	"fmt"
	"io"
	"os"
)

func MergeParts(output string, parts []*os.File) error {
	fmt.Println("下载完毕,开始合并")
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, part := range parts {
		part.Seek(0, 0) // 重置读取位置到文件开头
		_, err := io.Copy(outFile, part)
		if err != nil {
			return err
		}
		part.Close()
		os.Remove(part.Name()) // 删除临时文件
	}

	return nil
}
