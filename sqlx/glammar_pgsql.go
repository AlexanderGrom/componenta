package sqlx

import (
	"strconv"
	"strings"
)

type pgsqlGlammar struct {
	baseGlammar
	placeholders int
}

var _ glammar = (*pgsqlGlammar)(nil)

func init() {
	registerDriver("postgres", newPgsqlGlammar)
}

func newPgsqlGlammar() glammar {
	g := &pgsqlGlammar{}
	g.baseGlammar.glammar = g
	return g
}

// Обертка строки обратными кавычками
func (self *pgsqlGlammar) wrapp(value string) string {
	valtypes := strings.Split(value, "::")
	segments := strings.Split(valtypes[0], ".")
	wrapped := make([]string, len(segments))

	for k, v := range segments {
		if v == "*" {
			wrapped[k] = v
		} else {
			wrapped[k] = self.wrapQuote(v)
		}
	}

	if len(valtypes) > 1 {
		return strings.Join(wrapped, ".") + "::" + valtypes[1]
	}

	return strings.Join(wrapped, ".")
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

// Вставка Insert Returning id
func (self *pgsqlGlammar) compileReturning(b *Builder) string {
	if len(b.components.ReturnId) == 0 {
		return ""
	}
	return "RETURNING " + self.wrap("id")
}

// Вставка Insert ON CONFLICT DO NOTHING
func (self *pgsqlGlammar) compileOnConflictDoNothing(b *Builder) string {
	if len(b.components.OrIgnore) == 0 {
		return ""
	}
	return "ON CONFLICT DO NOTHING"
}
