package sqlx

import (
	"strings"
)

type glammar interface {
	wrapQuote(string) string
	parameter(...interface{}) string
	prepareRaw(interface{}) string
}

// Базовая граматика
type baseGlammar struct {
	glammar
}

func newBaseGlammar(g glammar) *baseGlammar {
	return &baseGlammar{g}
}

// Комбинация Select
func (self *baseGlammar) combineSelect(b *Builder) string {
	return combine(
		self.compileSelect(b),
		self.compileFrom(b),
		self.compileJoin(b),
		self.compileWhere(b),
		self.compileGroup(b),
		self.compileHaving(b),
		self.compileOrder(b),
		self.compileLimit(b),
		self.compileOffset(b),
	)
}

// Комбинация Delete
func (self *baseGlammar) combineDelete(b *Builder) string {
	return combine(
		self.compileDelete(b),
		self.compileFrom(b),
		self.compileWhere(b),
	)
}

// Комбинация Update
func (self *baseGlammar) combineUpdate(b *Builder) string {
	return combine(
		self.compileUpdate(b),
		self.compileSet(b),
		self.compileWhere(b),
	)
}

// Комбинация Insert
func (self *baseGlammar) combineInsert(b *Builder) string {
	return combine(
		self.compileInsert(b),
		self.compileInto(b),
		self.compileColumns(b),
		self.compileValues(b),
	)
}

// Компиляция Select
func (self *baseGlammar) compileSelect(b *Builder) string {
	buff := make([]interface{}, 0, len(b.components.Select)+len(b.components.Aggregate))
	buff = append(buff, self.selectFields(b)...)
	buff = append(buff, self.selectAggregates(b)...)
	return "SELECT " + self.wrap(buff...)
}

func (self *baseGlammar) selectFields(b *Builder) []interface{} {
	return b.components.Select
}

func (self *baseGlammar) selectAggregates(b *Builder) []interface{} {
	buff := make([]interface{}, len(b.components.Aggregate))
	for k, v := range b.components.Aggregate {
		buff[k] = Raw(strings.ToUpper(v.function) + "(" + self.wrap(v.column) + ") as " + self.wrap(v.alias))
	}
	return buff
}

// Компиляция From
func (self *baseGlammar) compileFrom(b *Builder) string {
	if len(b.components.From) == 0 {
		return ""
	}

	buff := make([]string, len(b.components.From))

	for k, v := range b.components.From {
		switch v.kind {
		case "str":
			buff[k] = self.formStr(v)
		case "sub":
			buff[k] = self.formSub(v)
		case "exp":
			buff[k] = self.formExp(v)
		}
	}

	return "FROM " + strings.Join(buff, ", ")
}

func (self *baseGlammar) formStr(f fromComponent) string {
	return self.wrap(f.table)
}

func (self *baseGlammar) formSub(f fromComponent) string {
	return combine("(", self.compile(f.builder), ")", "as", self.wrap(f.builder.table))
}

func (self *baseGlammar) formExp(f fromComponent) string {
	return self.prepareRaw(f.table)
}

// Компиляция Join
func (self *baseGlammar) compileJoin(b *Builder) string {
	if len(b.components.Join) == 0 {
		return ""
	}
	buff := make([]string, len(b.components.Join))
	for k, v := range b.components.Join {
		buff[k] = combine(v.kind, "JOIN", self.wrap(v.table), "ON (", self.wrap(v.column1), v.operator, self.wrap(v.column2), ")")
	}
	return strings.Join(buff, " ")
}

// Компиляция WHERE
func (self *baseGlammar) compileWhere(b *Builder) string {
	if len(b.components.Where) == 0 {
		return ""
	}

	buff := make([]string, len(b.components.Where))

	for k, v := range b.components.Where {
		var result string
		switch v.kind {
		case "base":
			result = self.whereBase(v)
		case "raw":
			result = self.whereRaw(v)
		case "group":
			result = self.whereGroup(v)
		case "between":
			result = self.whereBetween(v)
		case "notbetween":
			result = self.whereNotBetween(v)
		case "null":
			result = self.whereNull(v)
		case "notnull":
			result = self.whereNotNull(v)
		case "in":
			result = self.whereIn(v)
		case "notin":
			result = self.whereNotIn(v)
		case "insub":
			result = self.whereInSub(v)
		case "notinsub":
			result = self.whereNotInSub(v)
		}
		buff[k] = combine(v.boolean, result)
	}

	return "WHERE " + strings.Join(buff, " ")
}

func (self *baseGlammar) whereBase(w whereComponent) string {
	return combine(self.wrap(w.column), w.operator, self.parameter(w.value))
}

func (self *baseGlammar) whereGroup(w whereComponent) string {
	return "( " + self.compileWhere(w.builder)[6:] + " )"
}

func (self *baseGlammar) whereRaw(w whereComponent) string {
	return self.prepareRaw(w.value)
}

func (self *baseGlammar) whereIn(w whereComponent) string {
	return combine(self.wrap(w.column), "IN (", self.parameter(w.list...), ")")
}

func (self *baseGlammar) whereNotIn(w whereComponent) string {
	return combine(self.wrap(w.column), "NOT IN (", self.parameter(w.list...), ")")
}

func (self *baseGlammar) whereInSub(w whereComponent) string {
	return combine(self.wrap(w.column), "IN (", self.compile(w.builder), ")")
}

func (self *baseGlammar) whereNotInSub(w whereComponent) string {
	return combine(self.wrap(w.column), "NOT IN (", self.compile(w.builder), ")")
}

func (self *baseGlammar) whereBetween(w whereComponent) string {
	return combine(self.wrap(w.column), "BETWEEN", self.parameter(w.min), "AND", self.parameter(w.max))
}

func (self *baseGlammar) whereNotBetween(w whereComponent) string {
	return combine(self.wrap(w.column), "NOT BETWEEN", self.parameter(w.min), "AND", self.parameter(w.max))
}

func (self *baseGlammar) whereNull(w whereComponent) string {
	return self.wrap(w.column) + " IS NULL"
}

func (self *baseGlammar) whereNotNull(w whereComponent) string {
	return self.wrap(w.column) + " IS NOT NULL"
}

// Компиляция Having
func (self *baseGlammar) compileHaving(b *Builder) string {
	if len(b.components.Having) == 0 {
		return ""
	}

	buff := make([]string, len(b.components.Having))

	for k, v := range b.components.Having {
		var result string
		switch v.kind {
		case "base":
			result = self.havingBase(v)
		case "raw":
			result = self.havingRaw(v)
		case "group":
			result = self.havingGroup(v)
		}
		buff[k] = combine(v.boolean, result)
	}

	return "HAVING " + strings.Join(buff, " ")
}

func (self *baseGlammar) havingBase(h havingComponent) string {
	return combine(self.wrap(h.column), h.operator, self.parameter(h.value))
}

func (self *baseGlammar) havingRaw(h havingComponent) string {
	return self.prepareRaw(h.value)
}

func (self *baseGlammar) havingGroup(h havingComponent) string {
	return "( " + self.compileHaving(h.builder)[7:] + " )"
}

// Компиляция Group By
func (self *baseGlammar) compileGroup(b *Builder) string {
	if len(b.components.Group) == 0 {
		return ""
	}
	return "GROUP BY " + self.wrap(b.components.Group...)
}

// Компиляция Order By
func (self *baseGlammar) compileOrder(b *Builder) string {
	if len(b.components.Order) == 0 {
		return ""
	}

	buff := make([]string, len(b.components.Order))

	for k, v := range b.components.Order {
		buff[k] = self.wrap(v.column) + " " + v.direction
	}

	return "ORDER BY " + strings.Join(buff, ", ")
}

// Компиляция Limit
func (self *baseGlammar) compileLimit(b *Builder) string {
	if len(b.components.Limit) == 0 {
		return ""
	}
	return "LIMIT " + self.parameter(b.components.Limit...)
}

// Компиляция Offset
func (self *baseGlammar) compileOffset(b *Builder) string {
	if len(b.components.Offset) == 0 {
		return ""
	}
	return "OFFSET " + self.parameter(b.components.Offset...)
}

// Компиляция Delete
func (self *baseGlammar) compileDelete(b *Builder) string {
	return "DELETE"
}

// Компиляция Update
func (self *baseGlammar) compileUpdate(b *Builder) string {
	return "UPDATE " + self.wrap(b.table)
}

// Компиляция Set
func (self *baseGlammar) compileSet(b *Builder) string {
	data := Data(b.components.Set[0])
	buff := make([]string, 0, len(data))
	keys := data.Keys()
	for _, k := range keys {
		buff = append(buff, combine(self.wrap(k), "=", self.parameter(data[k])))
	}
	return "SET " + strings.Join(buff, ", ")
}

// Компиляция Insert
func (self *baseGlammar) compileInsert(b *Builder) string {
	return "INSERT"
}

// Компиляция Into
func (self *baseGlammar) compileInto(b *Builder) string {
	return "INTO " + self.wrap(b.components.Into...)
}

// Компиляция Columns
func (self *baseGlammar) compileColumns(b *Builder) string {
	return "( " + self.wrap(b.components.Columns...) + " )"
}

// Компиляция Values
func (self *baseGlammar) compileValues(b *Builder) string {
	buff := make([]string, len(b.components.Values))
	for k, v := range b.components.Values {
		buff[k] = "( " + self.parameter(v...) + " )"
	}
	return "VALUES " + strings.Join(buff, ", ")
}

// Компилируем Builder
func (self *baseGlammar) compile(b *Builder) string {
	result := ""
	switch b.kind {
	case "select":
		result = self.combineSelect(b)
	case "insert":
		result = self.combineInsert(b)
	case "update":
		result = self.combineUpdate(b)
	case "delete":
		result = self.combineDelete(b)
	}
	return result
}

// Обертывание значений обратными кавычками
func (self *baseGlammar) wrap(values ...interface{}) string {
	buf := make([]string, len(values))
	for k, v := range values {
		if exp, ok := v.(Expression); ok {
			buf[k] = self.prepareRaw(exp)
			continue
		}

		str := toString(v)
		segments := strings.Fields(str)
		lenght := len(segments)

		switch {
		case lenght == 2:
			buf[k] = self.wrapp(segments[0]) + " as " + self.wrapQuote(segments[1])
		case lenght == 3:
			buf[k] = self.wrapp(segments[0]) + " as " + self.wrapQuote(segments[2])
		default:
			buf[k] = self.wrapp(str)
		}
	}
	return strings.Join(buf, ", ")
}

// Обертка строки обратными кавычками
func (self *baseGlammar) wrapp(value string) string {
	segments := strings.Split(value, ".")
	wrapped := make([]string, len(segments))

	for k, v := range segments {
		if v == "*" {
			wrapped[k] = v
		} else {
			wrapped[k] = self.wrapQuote(v)
		}
	}

	return strings.Join(wrapped, ".")
}
