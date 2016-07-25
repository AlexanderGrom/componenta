package sqlx

import (
	"strconv"
	"strings"
)

type pgsqlGlammar struct {
	placeholders int
}

var _ glammar = (*pgsqlGlammar)(nil)

func newPgsqlGlammar() glammar {
	return &pgsqlGlammar{}
}

func (self *pgsqlGlammar) wrapQuote(v string) string {
	return `"` + v + `"`
}

func (self *pgsqlGlammar) parameter(p ...interface{}) string {
	params := make([]string, len(p))
	for k, v := range p {
		if exp, ok := v.(Expression); ok {
			params[k] = exp.String()
		} else {
			self.placeholders++
			params[k] = "$" + strconv.Itoa(self.placeholders)
		}
	}
	return strings.Join(params, ", ")
}
