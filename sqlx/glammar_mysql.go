package sqlx

import (
	"strings"
)

type mysqlGlammar struct {
}

var _ glammar = (*mysqlGlammar)(nil)

func newMysqlGlammar() glammar {
	return &mysqlGlammar{}
}

func (self *mysqlGlammar) wrapQuote(v string) string {
	return "`" + v + "`"
}

func (self *mysqlGlammar) parameter(p ...interface{}) string {
	params := make([]string, len(p))
	for k, v := range p {
		if exp, ok := v.(Expression); ok {
			params[k] = exp.String()
		} else {
			params[k] = "?"
		}
	}
	return strings.Join(params, ", ")
}
