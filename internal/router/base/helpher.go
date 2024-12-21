package base

type Pager struct {
	// 页码
	Page int `json:"page,omitempty"`
	// 每页数量
	PageSize int `json:"page_size,omitempty"`
	// 总行数
	TotalRows int `json:"total_rows,omitempty"`
}

type ResponseCommonStruct struct {
	Msg   string      `json:"msg,omitempty"`
	Code  int         `json:"code,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Pager Pager       `json:"pager,omitempty"`
	Count int         `json:"count,omitempty"`
}
