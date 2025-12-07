package pdf

import (
	"context"
	"fmt"
	pdfParser "github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/components/tool"
	einoutils "github.com/cloudwego/eino/components/tool/utils"
	"log"
	"strings"
)

func NewPdfToolWithEino(name string) tool.InvokableTool {
	// 用于告诉模型如何/何时/为什么使用这个工具
	// 可以在描述中包含少量示例
	toolDesc := `
	将本地PDF转换为纯文本
	需传入本地PDF的绝对路径
`
	pdfTool, err := einoutils.InferTool(name, toolDesc, convertPdfToText)
	if err != nil {
		log.Fatalf("NewPdfToolWithEino failed, err: %v", err)
	}
	log.Println("使用eino提供的 pdf 解析器, 初始化完成, 工具名称: ", name)

	return pdfTool
}

func convertPdfToText(ctx context.Context, req *ParsePdfRequest) (*ParsePdfResponse, error) {
	result := &ParsePdfResponse{
		Success: false,
		Meta:    genMeta(req),
	}

	file, err := validateAndOpenPdf(req)
	if err != nil {
		result.ErrorMsg = err.Error()
		return result, nil
	}
	defer file.Close()

	// 按大模型传入的参数决定是否分页
	einoPdfParser, err := pdfParser.NewPDFParser(ctx, &pdfParser.Config{ToPages: req.ToPages})
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("初始化eino pdf解析器失败: %v", err)
		return result, nil
	}

	docs, err := einoPdfParser.Parse(ctx, file,
		parser.WithURI(req.FilePath),
		parser.WithExtraMeta(result.Meta),
	)

	if err != nil {
		result.ErrorMsg = fmt.Sprintf("eino pdf解析器解析文件失败: %v", err)
		return result, nil
	}

	result.Success = true
	result.TotalPages = len(docs)

	if req.ToPages {
		// 分页
		// 分页模式：按页码整理文本（用索引+1作为页码，可靠无依赖）
		pages := make([]PdfPageText, 0, len(docs))
		for idx, doc := range docs {
			pages = append(pages, PdfPageText{
				Page:    idx + 1,
				Content: doc.Content,
			})
		}
		result.Pages = pages
	} else {
		// 不分页
		var sb strings.Builder
		for _, doc := range docs {
			sb.WriteString(doc.Content)
			sb.WriteString("\n")
		}
		result.Content = sb.String()
	}

	return result, nil
}
