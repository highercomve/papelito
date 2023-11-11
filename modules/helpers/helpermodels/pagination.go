package helpermodels

// Pagination paginable response
type Pagination struct {
	PageSizes   []int  `json:"-"`
	ServiceURL  string `json:"resource"`
	PageSize    int    `json:"page_size"`
	PageOffset  int    `json:"page_offset"`
	CurrentPage int    `json:"current_page"`
	Total       int    `json:"total"`
	Next        string `json:"next"`
	Prev        string `json:"prev"`
}
