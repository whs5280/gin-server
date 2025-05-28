package pagination

const PageTypePage = "page"
const PageTypeIndex = "index"

type ArgPagination struct {
	Page       int32  `form:"page,default=1"`
	PageSize   int32  `form:"page_size,default=20"`
	OrderIndex string `form:"order_index"`
	PageType   string `form:"page_type,default=page"` //分页方式，page:传统分页，index:OrderIndex方式
}

type Pagination struct {
	IsOver bool   `json:"is_over"`
	Params string `json:"params"`
}
