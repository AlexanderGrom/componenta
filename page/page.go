package page

import (
	"math"
	"strconv"
)

// Параметры пагинации
type Param struct {
	TotalItem   int
	ViewItem    int
	ViewPage    int
	CurrentPage int
	CurrentURI  string
	LimitPage   int
}

// Результат
type Result struct {
	List       []map[string]string
	Start      int
	Step       int
	Total      int
	Prev       int
	Current    int
	Next       int
	Exists     bool
	PrevURI    string
	CurrentURI string
	NextURI    string
	FirstURI   string
	LastURI    string
}

// Page
type Page struct {
	result *Result
}

// Конструктор
// p := New(&Param{
//        TotalItem:   200,
//        ViewItem:    10,
//        ViewPage:    5,
//        CurrentPage: 1,
//        CurrentURI:  "/test/best/?one=1&two=2",
//    }).Path("page")
func New(param *Param) *Page {
	currentPage := param.CurrentPage
	if currentPage <= 0 {
		currentPage = 1
	}

	existsPage := true
	if param.CurrentPage <= 0 {
		existsPage = false
	}

	totalItem := param.TotalItem
	if totalItem < 0 {
		totalItem = 0
	}

	viewItem := param.ViewItem
	if viewItem < 0 {
		viewItem = 0
	}

	viewPage := param.ViewPage
	if viewPage < 0 {
		viewPage = 0
	}

	start := (currentPage * viewItem) - viewItem
	currentURI := param.CurrentURI

	if param.LimitPage > 0 {
		limitPage := param.ViewItem * param.LimitPage
		if totalItem > limitPage {
			totalItem = limitPage
		}
	}

	totalPage := int(math.Ceil(float64(totalItem) / float64(viewItem)))
	if totalPage == 0 {
		totalPage = 1
	}

	listPage := make([]map[string]string, 0)

	if viewPage != 0 {
		pageParts := int(math.Floor(float64(viewPage) / 2.0))

		offset := 1
		if currentPage > (pageParts + 1) {
			offset = currentPage - pageParts
		}
		if currentPage > (totalPage - pageParts) {
			offset = totalPage - viewPage + 1
		}
		if offset <= 0 {
			offset = 1
		}

		iter := viewPage
		if currentPage > pageParts {
			iter = currentPage + pageParts
		}
		if iter > totalPage {
			iter = totalPage
		}

		for i := offset; i <= iter; i++ {
			listPage = append(listPage, map[string]string{
				"num": strconv.Itoa(i),
				"url": "",
			})
		}
	}

	prevPage := currentPage - 1
	nextPage := currentPage + 1

	return &Page{
		result: &Result{
			List:       listPage,
			Start:      start,
			Step:       viewItem,
			Total:      totalPage,
			Prev:       prevPage,
			Current:    currentPage,
			Next:       nextPage,
			Exists:     existsPage,
			CurrentURI: currentURI,
		},
	}
}
