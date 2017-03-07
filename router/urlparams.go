package router

// Хранилице для параметров URL
type URLParams map[string]string

func NewURLParams() URLParams {
	return URLParams{}
}

func (self URLParams) Get(key string) string {
	return self[key]
}

func (self URLParams) Exists(key string) bool {
	_, ok := self[key]
	return ok
}
