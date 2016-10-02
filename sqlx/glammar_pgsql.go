package sqlx

import (
	"strconv"
	"strings"
)

type pgsqlGlammar struct {
	placeholders int
}

var _ glammar = (*pgsqlGlammar)(nil)

func init() {
	registerDriver("pgsql", newPgsqlGlammar)
}

func newPgsqlGlammar() glammar {
	return &pgsqlGlammar{}
}

func (self *pgsqlGlammar) wrapQuote(v string) string {
	return `"` + v + `"`
}

func (self *pgsqlGlammar) placeholder() string {
	return "$" + strconv.Itoa(self.placeholders)
}

func (self *pgsqlGlammar) parameter(p ...interface{}) string {
	params := make([]string, len(p))
	for k, v := range p {
		if exp, ok := v.(Expression); ok {
			params[k] = self.prepareRaw(exp)
		} else {
			self.placeholders++
			params[k] = self.placeholder()
		}
	}
	return strings.Join(params, ", ")
}

func (self *pgsqlGlammar) prepareRaw(p interface{}) string {
	if exp, ok := p.(Expression); ok {
		str := exp.String()
		for range exp.Data() {
			self.placeholders++
			str = strings.Replace(str, "?", self.placeholder(), 1)
		}
		return str
	}
	return toString(p)
}
