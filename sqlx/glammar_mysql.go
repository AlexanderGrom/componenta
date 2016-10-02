package sqlx

import (
	"strings"
)

type mysqlGlammar struct {
}

var _ glammar = (*mysqlGlammar)(nil)

func init() {
	registerDriver("mysql", newMysqlGlammar)
}

func newMysqlGlammar() glammar {
	return &mysqlGlammar{}
}

func (self *mysqlGlammar) wrapQuote(v string) string {
	return "`" + v + "`"
}

func (self *mysqlGlammar) placeholder() string {
	return "?"
}

func (self *mysqlGlammar) parameter(p ...interface{}) string {
	params := make([]string, len(p))
	for k, v := range p {
		if exp, ok := v.(Expression); ok {
			params[k] = self.prepareRaw(exp)
		} else {
			params[k] = self.placeholder()
		}
	}
	return strings.Join(params, ", ")
}

func (self *mysqlGlammar) prepareRaw(p interface{}) string {
	return toString(p)
}
