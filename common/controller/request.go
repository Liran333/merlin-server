package controller

// reqToListUserModels
type CommonListRequest struct {
	SortBy       string `form:"sort_by"`
	Count        bool   `form:"count"`
	PageNum      int    `form:"page_num"`
	CountPerPage int    `form:"count_per_page"`
}
