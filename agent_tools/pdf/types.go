package pdf

type ParsePdfRequest struct {
	FilePath string `json:"filePath" jsonschema:"description=本地PDF文件的绝对路径"`
	ToPages  bool   `json:"toPages" jsonschema:"description=是否按页分割文本,true=分页输出"`
}

type ParsePdfResponse struct {
	Success    bool          `json:"success" jsonschema:"description=是否解析成功"`
	Content    string        `json:"content" jsonschema:"description=解析出的文本内容"`
	Pages      []PdfPageText `json:"pages" jsonschema:"description=按页分割的文本内容(仅当toPages为true时返回)"`
	TotalPages int           `json:"totalPages" jsonschema:"description=总页数"`
	ErrorMsg   string        `json:"errorMsg,omitempty" jsonschema:"description=解析失败时的错误信息"`
	Meta       ParsePdfMeta  `json:"meta,omitempty" jsonschema:"description=解析时的元数据"`
}

type ParsePdfMeta = map[string]any

type PdfPageText struct {
	Page    int    `json:"page" jsonschema:"description=页码"`
	Content string `json:"content" jsonschema:"description=该页的文本内容"`
}
