package lib

type Page struct {
	currentPage int64 // 当前页
	pageSize    int64 // 每页几行
	totalPage   int64 // 总页数
	totalCount  int64 // 总行数
	rangeCount  int64 //区间页索引个数
}

func NewPagerR(currentPage, pageSize, totalCount, rangeCount int64) *Page {
	p := &Page{
		currentPage: currentPage,
		pageSize:    pageSize,
		totalCount:  totalCount,
		rangeCount:  rangeCount,
	}
	p.setTotalPage()
	return p
}

func NewPager(currentPage, pageSize, totalCount int64) *Page {
	p := &Page{
		currentPage: currentPage,
		pageSize:    pageSize,
		totalCount:  totalCount,
		rangeCount:  int64(10),
	}
	p.setTotalPage()
	return p
}

func (p *Page) setTotalPage() {
	totalPage := p.TotalCount() / p.PageSize()
	if p.TotalCount()%p.PageSize() == 0 {
		p.totalPage = totalPage
	} else {
		p.totalPage = totalPage + 1
	}
}

func (p Page) CurrentPage() int64 {
	if p.currentPage <= 0 {
		p.currentPage = 1
	}
	return p.currentPage
}

func (p Page) PageSize() int64 {
	return p.pageSize
}

func (p Page) TotalPage() int64 {
	return p.totalPage
}

func (p Page) TotalCount() int64 {
	return p.totalCount
}

func (p Page) HasPrev() bool {
	if p.currentPage <= p.totalPage && p.currentPage > 1 {
		return true
	} else {
		return false
	}
}

func (p Page) HasNext() bool {
	if p.currentPage < p.totalPage {
		return true
	}
	return false
}

func (p Page) PrevPage() int64 {
	if p.currentPage-1 > 0 {
		return p.currentPage - 1
	} else {
		return 1
	}
}

func (p Page) NextPage() int64 {
	if p.currentPage+1 <= p.totalPage {
		return p.currentPage + 1
	}
	return p.currentPage
}

func (p Page) RangePage() []int64 {
	ret := make([]int64, 0)
	if p.totalPage > 0 {
		index := 0
		if p.rangeCount >= p.totalPage {
			for i := int64(1); i <= p.totalPage; i++ {
				ret = append(ret, i)
				index++
			}
		} else {
			if p.currentPage+p.rangeCount-1 <= p.totalPage {
				for i := p.currentPage; i <= p.currentPage+p.rangeCount-1; i++ {
					ret = append(ret, i)
					index++
				}
			} else {
				start := p.currentPage - (p.rangeCount - (p.totalPage - p.currentPage + 1))
				if start <= 0 {
					start = 1
				}
				for i := start; i <= p.totalPage; i++ {
					ret = append(ret, i)
					index++
				}
			}
		}
	}
	return ret
}

/**
func main() {
	p := NewPager(1, 20, 10000)
	fmt.Printf("上一页：%d \n", p.PrevPage())
	fmt.Printf("下一页：%d \n", p.NextPage())
	fmt.Printf("当前页：%d \n", p.CurrentPage())
	fmt.Printf("每页记录数：%d \n", p.PageSize())
	fmt.Printf("总记录数：%d \n", p.TotalCount())
	fmt.Printf("总页数：%d \n", p.TotalPage())
	//fmt.Println(p.RangePage())
}
**/
