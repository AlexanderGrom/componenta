package sqlx

import (
	"strings"
)

type sqliteGlammar struct {
}

var _ glammar = (*sqliteGlammar)(nil)

func init() {
	registerDriver("mysql", newSqliteGlammar)
}

func newSqliteGlammar() glammar {
	return &sqliteGlammar{}
}

func (self *sqliteGlammar) wrapQuote(v string) string {
	return "`" + v + "`"
}

func (self *sqliteGlammar) placeholder() string {
	return "?"
}

func (self *sqliteGlammar) parameter(p ...interface{}) string {
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

func (self *sqliteGlammar) prepareRaw(p interface{}) string {
	return toString(p)
}
