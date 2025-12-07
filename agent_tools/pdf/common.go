package pdf

import (
	"fmt"
	"os"
	"time"
)

func genMeta(req *ParsePdfRequest) ParsePdfMeta {
	meta := make(ParsePdfMeta)
	meta["filePath"] = req.FilePath
	meta["toPages"] = req.ToPages
	meta["parseTime"] = time.Now().Format("2006-01-02 15:04:05")
	return meta
}

// 校验并打开PDF文件
func validateAndOpenPdf(req *ParsePdfRequest) (*os.File, error) {
	if req.FilePath == "" {
		return nil, fmt.Errorf("必须传入 pdf 文件的绝对路径")
	}
	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("打开PDF文件失败：%v（请检查路径是否正确、文件是否存在）", err)
	}
	return file, err
}
