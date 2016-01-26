
## Componenta / Page

Набросок пакета расчета пагинации

```go
package main

import (
    "github.com/AlexanderGrom/componenta/page"
)

func main() {
    p := page.New(&page.Param{
        TotalItem:   200,   // Общее кол-во элементов
        ViewItem:    10,    // Кол-во элементов которые нужно показать на одной странице
        ViewPage:    5,     // Кол-во показываемых страниц к пагинации
        CurrentPage: 1,     // Номер текущей страницы
        CurrentURI:  "/test/?one=1&two=2", // Текущий URI страницы
    }).Path("page") // page это префикс сигмента с номером страницы /test/page5/?one=1&two=2
    
    fmt.Printf("%v\n", p)
    
    /*
    //Path("page") вернет структуру *page.Result c данными о текущей странице
    
    type Result struct {
        List       []map[string]string // элементы пагинации number, url.
        Start      int // число записей, которые нужно пропустить при выборке.
        Step       int // число записей, которые нужно выбрать для текущей странице.
        Total      int // общее кол-во доступных страниц.
        Prev       int // номер предыдущей страницы.
        Current    int // номер текущей страницы.
        Next       int // номер следующей страницы.
        Exists     bool // если Param.CurrentPage это положительное число.
        PrevURI    string // uri предыдущей страницы.
        CurrentURI string // uri текущей страницы.
        NextURI    string // uri следующей страницы.
        FirstURI   string // uri первой страницы.
        LastURI    string // uri последней страницы.
    }
    */
    
    // Так же можно размещать номер страницы в разделе Query
    // p := page.New(&Param{/*...*/}).Query("page") // page это параметр со значеним номера страницы /test?one=1&two=2&page=5
}
```