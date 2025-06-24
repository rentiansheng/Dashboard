package define

/***************************
    @author: tiansheng.ren
    @date: 2025/6/9
    @desc:

***************************/

type Page struct {
	PageNum  int64 `json:"page_num"  query:"page_num"`
	PageSize int64 `json:"page_size"  query:"page_size" validate:"max=200"`
}

func NewPage(pageNum, pageSize int64) *Page {
	p := &Page{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
	p.Default()
	return p
}

func NewPageWithDefault() *Page {
	p := &Page{}
	p.Default()
	return p
}

func NewFirstPage(pageSize int64) *Page {
	return NewPage(1, pageSize)
}

func NewFirstPage1k() *Page {
	return NewFirstPage(1000)
}

func (p *Page) Default() {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
}

func (p Page) Offset() int64 {
	if p.PageNum <= 0 {
		return 0
	}
	return p.Limit() * (p.PageNum - 1)
}

// Limit 如果PageSize为0，则默认返回1000条
func (p Page) Limit() int64 {
	if p.PageSize <= 0 {
		return 1000
	}
	return p.PageSize
}

func (p *Page) Next() {
	p.PageNum += 1
}
