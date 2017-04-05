package sqlx

import (
	"strings"
)

// Интерфейс, описываем тут те методы
// для которых хотим иметь возможность переопределения
type glammar interface {
	wrap(...interface{}) string
	wrapp(value string) string
	wrapQuote(string) string
	parameter(...interface{}) string
	prepareRaw(interface{}) string
	combineSelect(*Builder) string
	combineInsert(*Builder) string
	combineUpdate(*Builder) string
	combineDelete(*Builder) string
	compile(*Builder) string
	compileSelect(*Builder) string
	compileFrom(*Builder) string
	compileJoin(*Builder) string
	compileWhere(*Builder) string
	compileGroup(*Builder) string
	compileHaving(*Builder) string
	compileOrder(*Builder) string
	compileLimit(*Builder) string
	compileOffset(*Builder) string
	compileDelete(*Builder) string
	compileUpdate(*Builder) string
	compileSet(*Builder) string
	compileInsert(*Builder) string
	compileInto(*Builder) string
	compileColumns(*Builder) string
	compileValues(*Builder) string
	compileReturning(*Builder) string
}

// Базовая граматика
type baseGlammar struct {
	glammar
}

// Комбинация Select
func (self *baseGlammar) combineSelect(b *Builder) string {
	return combine(
		self.glammar.compileSelect(b),
		self.glammar.compileFrom(b),
		self.glammar.compileJoin(b),
		self.glammar.compileWhere(b),
		self.glammar.compileGroup(b),
		self.glammar.compileHaving(b),
		self.glammar.compileOrder(b),
		self.glammar.compileLimit(b),
		self.glammar.compileOffset(b),
	)
}

// Комбинация Delete
func (self *baseGlammar) combineDelete(b *Builder) string {
	return combine(
		self.glammar.compileDelete(b),
		self.glammar.compileFrom(b),
		self.glammar.compileWhere(b),
	)
}

// Комбинация Update
func (self *baseGlammar) combineUpdate(b *Builder) string {
	return combine(
		self.glammar.compileUpdate(b),
		self.glammar.compileSet(b),
		self.glammar.compileWhere(b),
	)
}

// Комбинация Insert
func (self *baseGlammar) combineInsert(b *Builder) string {
	return combine(
		self.glammar.compileInsert(b),
		self.glammar.compileInto(b),
		self.glammar.compileColumns(b),
		self.glammar.compileValues(b),
		self.glammar.compileReturning(b),
	)
}

// Компиляция Select
func (self *baseGlammar) compileSelect(b *Builder) string {
	buff := make([]interface{}, 0, len(b.components.Select)+len(b.components.Aggregate))
	buff = append(buff, self.selectFields(b)...)
	buff = append(buff, self.selectAggregates(b)...)
	return "SELECT " + self.glammar.wrap(buff...)
}

func (self *baseGlammar) selectFields(b *Builder) []interface{} {
	return b.components.Select
}

func (self *baseGlammar) selectAggregates(b *Builder) []interface{} {
	buff := make([]interface{}, len(b.components.Aggregate))
	for k, v := range b.components.Aggregate {
		buff[k] = Raw(strings.ToUpper(v.function) + "(" + self.glammar.wrap(v.column) + ") as " + self.glammar.wrap(v.alias))
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
	return self.glammar.wrap(f.table)
}

func (self *baseGlammar) formSub(f fromComponent) string {
	return combine("(", self.glammar.compile(f.builder), ")", "as", self.glammar.wrap(f.builder.table))
}

func (self *baseGlammar) formExp(f fromComponent) string {
	return self.glammar.prepareRaw(f.table)
}

// Компиляция Join
func (self *baseGlammar) compileJoin(b *Builder) string {
	if len(b.components.Join) == 0 {
		return ""
	}

	buff1 := make([]string, len(b.components.Join))
	for i, join := range b.components.Join {
		buff2 := make([]string, len(join.conditions))
		for j, v := range join.conditions {
			var result string
			switch v.kind {
			case "on":
				result = combine(self.glammar.wrap(v.column), v.operator, self.glammar.wrap(v.value))
			case "where":
				result = combine(self.glammar.wrap(v.column), v.operator, self.glammar.parameter(v.value))
			}
			buff2[j] = combine(v.boolean, result)
		}

		buff1[i] = combine(join.kind, "JOIN", self.glammar.wrap(join.table), "ON (", strings.Join(buff2, " "), ")")
	}

	return strings.Join(buff1, " ")
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
	return combine(self.glammar.wrap(w.column), w.operator, self.glammar.parameter(w.value))
}

func (self *baseGlammar) whereGroup(w whereComponent) string {
	return "( " + self.glammar.compileWhere(w.builder)[6:] + " )"
}

func (self *baseGlammar) whereRaw(w whereComponent) string {
	return self.glammar.prepareRaw(w.value)
}

func (self *baseGlammar) whereIn(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "IN (", self.glammar.parameter(w.list...), ")")
}

func (self *baseGlammar) whereNotIn(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "NOT IN (", self.glammar.parameter(w.list...), ")")
}

func (self *baseGlammar) whereInSub(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "IN (", self.glammar.compile(w.builder), ")")
}

func (self *baseGlammar) whereNotInSub(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "NOT IN (", self.glammar.compile(w.builder), ")")
}

func (self *baseGlammar) whereBetween(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "BETWEEN", self.glammar.parameter(w.min), "AND", self.glammar.parameter(w.max))
}

func (self *baseGlammar) whereNotBetween(w whereComponent) string {
	return combine(self.glammar.wrap(w.column), "NOT BETWEEN", self.glammar.parameter(w.min), "AND", self.glammar.parameter(w.max))
}

func (self *baseGlammar) whereNull(w whereComponent) string {
	return self.glammar.wrap(w.column) + " IS NULL"
}

func (self *baseGlammar) whereNotNull(w whereComponent) string {
	return self.glammar.wrap(w.column) + " IS NOT NULL"
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
	return combine(self.glammar.wrap(h.column), h.operator, self.glammar.parameter(h.value))
}

func (self *baseGlammar) havingRaw(h havingComponent) string {
	return self.glammar.prepareRaw(h.value)
}

func (self *baseGlammar) havingGroup(h havingComponent) string {
	return "( " + self.glammar.compileHaving(h.builder)[7:] + " )"
}

// Компиляция Group By
func (self *baseGlammar) compileGroup(b *Builder) string {
	if len(b.components.Group) == 0 {
		return ""
	}
	return "GROUP BY " + self.glammar.wrap(b.components.Group...)
}

// Компиляция Order By
func (self *baseGlammar) compileOrder(b *Builder) string {
	if len(b.components.Order) == 0 {
		return ""
	}

	buff := make([]string, len(b.components.Order))

	for k, v := range b.components.Order {
		buff[k] = self.glammar.wrap(v.column) + " " + v.direction
	}

	return "ORDER BY " + strings.Join(buff, ", ")
}

// Компиляция Limit
func (self *baseGlammar) compileLimit(b *Builder) string {
	if len(b.components.Limit) == 0 {
		return ""
	}
	return "LIMIT " + self.glammar.parameter(b.components.Limit...)
}

// Компиляция Offset
func (self *baseGlammar) compileOffset(b *Builder) string {
	if len(b.components.Offset) == 0 {
		return ""
	}
	return "OFFSET " + self.glammar.parameter(b.components.Offset...)
}

// Компиляция Delete
func (self *baseGlammar) compileDelete(b *Builder) string {
	return "DELETE"
}

// Компиляция Update
func (self *baseGlammar) compileUpdate(b *Builder) string {
	return "UPDATE " + self.glammar.wrap(b.table)
}

// Компиляция Set
func (self *baseGlammar) compileSet(b *Builder) string {
	data := Data(b.components.Set[0])
	buff := make([]string, 0, len(data))
	keys := data.Keys()
	for _, k := range keys {
		buff = append(buff, combine(self.glammar.wrap(k), "=", self.glammar.parameter(data[k])))
	}
	return "SET " + strings.Join(buff, ", ")
}

// Компиляция Insert
func (self *baseGlammar) compileInsert(b *Builder) string {
	return "INSERT"
}

// Компиляция Into
func (self *baseGlammar) compileInto(b *Builder) string {
	return "INTO " + self.glammar.wrap(b.components.Into...)
}

// Компиляция Columns
func (self *baseGlammar) compileColumns(b *Builder) string {
	return "( " + self.glammar.wrap(b.components.Columns...) + " )"
}

// Компиляция Values
func (self *baseGlammar) compileValues(b *Builder) string {
	buff := make([]string, len(b.components.Values))
	for k, v := range b.components.Values {
		buff[k] = "( " + self.glammar.parameter(v...) + " )"
	}
	return "VALUES " + strings.Join(buff, ", ")
}

// Вставка Insert Returning id (заглушка)
func (self *baseGlammar) compileReturning(b *Builder) string {
	return ""
}

// Компилируем Builder
func (self *baseGlammar) compile(b *Builder) string {
	result := ""
	switch b.kind {
	case "select":
		result = self.glammar.combineSelect(b)
	case "insert":
		result = self.glammar.combineInsert(b)
	case "update":
		result = self.glammar.combineUpdate(b)
	case "delete":
		result = self.glammar.combineDelete(b)
	}
	return result
}

// Обертывание значений обратными кавычками
func (self *baseGlammar) wrap(values ...interface{}) string {
	buf := make([]string, len(values))
	for k, v := range values {
		if exp, ok := v.(Expression); ok {
			buf[k] = self.glammar.prepareRaw(exp)
			continue
		}

		str := toString(v)
		segments := strings.Fields(str)
		lenght := len(segments)

		switch {
		case lenght == 2:
			buf[k] = self.glammar.wrapp(segments[0]) + " as " + self.glammar.wrapQuote(segments[1])
		case lenght == 3:
			buf[k] = self.glammar.wrapp(segments[0]) + " as " + self.glammar.wrapQuote(segments[2])
		default:
			buf[k] = self.glammar.wrapp(str)
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
			wrapped[k] = self.glammar.wrapQuote(v)
		}
	}

	return strings.Join(wrapped, ".")
}

func (self *baseGlammar) wrapQuote(v string) string {
	return "`" + v + "`"
}

func (self *baseGlammar) placeholder() string {
	return "?"
}

func (self *baseGlammar) parameter(p ...interface{}) string {
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

func (self *baseGlammar) prepareRaw(p interface{}) string {
	return toString(p)
}
