package pagination

type Pagination struct {
	IsOver   bool `json:"is_over"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
}

func MakePagination[T any](list []T, page int, pageSize int) (pagination *Pagination) {
	return &Pagination{
		IsOver:   len(list) < pageSize,
		Page:     page,
		PageSize: pageSize,
	}
}
