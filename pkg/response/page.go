package response

// PageResult 通用分页返回结构
type PageResult struct {
	Total   int64       `json:"total"`   // 总记录数
	Records interface{} `json:"records"` // 数据集合
}

// NewPageResult 构造函数
func NewPageResult(total int64, records interface{}) *PageResult {
	return &PageResult{
		Total:   total,
		Records: records,
	}
}
