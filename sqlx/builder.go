package sqlx

import (
	"sort"
)

type glammarFunc func() glammar

var driver glammarFunc
var drivers = map[string]glammarFunc{
	"mysql": newMysqlGlammar,
	"pgsql": newPgsqlGlammar,
}

type List []interface{}
type Data map[string]interface{}

func (self Data) Keys() []string {
	keys := make([]string, 0, len(self))
	for k := range self {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (self Data) Values() []interface{} {
	keys := self.Keys()
	values := make([]interface{}, 0, len(self))
	for _, k := range keys {
		values = append(values, self[k])
	}
	return values
}

// Драйвер грамматики, который будет использован для построения запроса
// Параметр name может принимать значения: mysql, pgsql
func Driver(name string) {
	glammar, ok := drivers[name]
	if !ok {
		panic("sqlx: driver '" + name + "' not found")
	}
	driver = glammar
}

// Helper для создания строителя
func Table(table string) *Builder {
	return NewBuilder().From(table)
}

// Kарта значений для плейсехолдеров
var bindings = []string{"values", "set", "where", "having", "limit", "offset"}

// Карта значений для плейсехолдеров в зависимости от типа запроса
var bindingsMap = map[string][]string{
	"select": []string{"where", "having", "limit", "offset"},
	"insert": []string{"values"},
	"update": []string{"set", "where"},
	"delete": []string{"where"},
}

// Создаем карту значений для плейсехолдеров
func newBindings() map[string][]interface{} {
	b := make(map[string][]interface{})
	for _, v := range bindings {
		b[v] = []interface{}{}
	}
	return b
}

type Builder struct {
	kind       string
	table      string
	components *components
	bindings   map[string][]interface{}
}

func NewBuilder() *Builder {
	return &Builder{
		components: newComponents(),
		bindings:   newBindings(),
	}
}

func (self *Builder) Select(p ...interface{}) *Builder {
	if self.kind != "" && self.kind != "select" {
		return self
	}
	self.kind = "select"
	self.components.Select = append(self.components.Select, p...)
	return self
}

func (self *Builder) SelectRaw(exp ...string) *Builder {
	buff := make([]interface{}, len(exp))
	for k, v := range exp {
		buff[k] = Raw(v)
	}
	return self.Select(buff...)
}

func (self *Builder) From(name string) *Builder {
	self.table = name
	self.components.From = []interface{}{name}
	return self
}

func (self *Builder) Join(table, column1, operator, column2 string) *Builder {
	self.join(table, column1, operator, column2, "INNER")
	return self
}

func (self *Builder) LeftJoin(table, column1, operator, column2 string) *Builder {
	self.join(table, column1, operator, column2, "LEFT")
	return self
}

func (self *Builder) join(table, column1, operator, column2, kind string) {
	self.components.Join = append(self.components.Join, joinComponent{
		kind:     kind,
		table:    table,
		column1:  column1,
		operator: operator,
		column2:  column2,
	})
}

func (self *Builder) Where(column string, operator string, value interface{}) *Builder {
	self.where(column, operator, value, "AND")
	return self
}

func (self *Builder) OrWhere(column string, operator string, value interface{}) *Builder {
	self.where(column, operator, value, "OR")
	return self
}

func (self *Builder) where(column string, operator string, value interface{}, boolean string) {
	if !isOperator(operator) {
		panic("sqlx: such a \"operator\" is not allowed")
	}
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	self.components.Where = append(self.components.Where, whereComponent{
		kind:     "base",
		column:   column,
		operator: operator,
		value:    value,
		boolean:  boolean,
	})
	self.bind("where", value)
}

func (self *Builder) WhereGroup(callback func(*Builder)) *Builder {
	self.whereGroup(callback, "AND")
	return self
}

func (self *Builder) OrWhereGroup(callback func(*Builder)) *Builder {
	self.whereGroup(callback, "OR")
	return self
}

func (self *Builder) whereGroup(callback func(*Builder), boolean interface{}) {
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	builder := Table(self.table)
	callback(builder)
	if len(builder.components.Where) > 0 {
		self.components.Where = append(self.components.Where, whereComponent{
			kind:    "group",
			builder: builder,
			boolean: boolean,
		})
		self.bind("where", builder.Data()...)
	}
}

func (self *Builder) WhereRaw(exp string) *Builder {
	self.whereRaw(exp, "AND")
	return self
}

func (self *Builder) OrWhereRaw(exp string) *Builder {
	self.whereRaw(exp, "OR")
	return self
}

func (self *Builder) whereRaw(exp string, boolean interface{}) {
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	self.components.Where = append(self.components.Where, whereComponent{
		kind:    "raw",
		value:   exp,
		boolean: boolean,
	})
}

func (self *Builder) WhereBetween(column string, min, max interface{}) *Builder {
	self.whereBetween(column, min, max, "AND", false)
	return self
}

func (self *Builder) OrWhereBetween(column string, min, max interface{}) *Builder {
	self.whereBetween(column, min, max, "OR", false)
	return self
}

func (self *Builder) WhereNotBetween(column string, min, max interface{}) *Builder {
	self.whereBetween(column, min, max, "AND", true)
	return self
}

func (self *Builder) OrWhereNotBetween(column string, min, max interface{}) *Builder {
	self.whereBetween(column, min, max, "OR", true)
	return self
}

func (self *Builder) whereBetween(column string, min, max interface{}, boolean string, not bool) {
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	kind := "between"
	if not {
		kind = "notbetween"
	}
	self.components.Where = append(self.components.Where, whereComponent{
		kind:    kind,
		column:  column,
		min:     min,
		max:     max,
		boolean: boolean,
	})
	self.bind("where", min, max)
}

func (self *Builder) WhereNull(column string) *Builder {
	self.whereNull(column, "AND", false)
	return self
}

func (self *Builder) OrWhereNull(column string) *Builder {
	self.whereNull(column, "OR", false)
	return self
}

func (self *Builder) WhereNotNull(column string) *Builder {
	self.whereNull(column, "AND", true)
	return self
}

func (self *Builder) OrWhereNotNull(column string) *Builder {
	self.whereNull(column, "OR", true)
	return self
}

func (self *Builder) whereNull(column string, boolean string, not bool) {
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	kind := "null"
	if not {
		kind = "notnull"
	}
	self.components.Where = append(self.components.Where, whereComponent{
		kind:    kind,
		column:  column,
		boolean: boolean,
	})
}

func (self *Builder) WhereIn(column string, values interface{}) *Builder {
	self.whereIn(column, values, "AND", false)
	return self
}

func (self *Builder) OrWhereIn(column string, values interface{}) *Builder {
	self.whereIn(column, values, "OR", false)
	return self
}

func (self *Builder) WhereNotIn(column string, values interface{}) *Builder {
	self.whereIn(column, values, "AND", true)
	return self
}

func (self *Builder) OrWhereNotIn(column string, values interface{}) *Builder {
	self.whereIn(column, values, "OR", true)
	return self
}

func (self *Builder) whereIn(column string, values interface{}, boolean string, not bool) {
	if len(self.components.Where) == 0 {
		boolean = ""
	}
	switch v := values.(type) {
	case List:
		self.whereInList(column, v, boolean, not)
	case func(*Builder):
		self.whereInSub(column, v, boolean, not)
	default:
		self.whereInList(column, List{values}, boolean, not)
	}
}

func (self *Builder) whereInList(column string, list List, boolean string, not bool) {
	kind := "in"
	if not {
		kind = "notin"
	}
	self.components.Where = append(self.components.Where, whereComponent{
		kind:    kind,
		column:  column,
		list:    list,
		boolean: boolean,
	})
	self.bind("where", list...)
}

func (self *Builder) whereInSub(column string, callback func(*Builder), boolean string, not bool) {
	kind := "insub"
	if not {
		kind = "notinsub"
	}
	builder := NewBuilder()
	callback(builder)
	self.components.Where = append(self.components.Where, whereComponent{
		kind:    kind,
		column:  column,
		builder: builder,
		boolean: boolean,
	})
	self.bind("where", builder.Data()...)
}

func (self *Builder) GroupBy(p ...interface{}) *Builder {
	self.components.Group = append(self.components.Group, p...)
	return self
}

func (self *Builder) Having(column string, operator string, value interface{}) *Builder {
	self.having(column, operator, value, "AND")
	return self
}

func (self *Builder) OrHaving(column string, operator string, value interface{}) *Builder {
	self.having(column, operator, value, "OR")
	return self
}

func (self *Builder) having(column string, operator string, value interface{}, boolean string) {
	if !isOperator(operator) {
		panic("sqlx: operator is not allowed")
	}
	if len(self.components.Having) == 0 {
		boolean = ""
	}
	self.components.Having = append(self.components.Having, havingComponent{
		kind:     "base",
		column:   column,
		operator: operator,
		value:    value,
		boolean:  boolean,
	})
	self.bind("having", value)
}

func (self *Builder) HavingGroup(callback func(*Builder)) *Builder {
	self.havingGroup(callback, "AND")
	return self
}

func (self *Builder) OrHavingGroup(callback func(*Builder)) *Builder {
	self.havingGroup(callback, "OR")
	return self
}

func (self *Builder) havingGroup(callback func(*Builder), boolean string) {
	if len(self.components.Having) == 0 {
		boolean = ""
	}
	builder := Table(self.table)
	callback(builder)
	if len(builder.components.Having) > 0 {
		self.components.Having = append(self.components.Having, havingComponent{
			kind:    "group",
			builder: builder,
			boolean: boolean,
		})
		self.bind("having", builder.Data()...)
	}
}

func (self *Builder) HavingRaw(exp string) *Builder {
	self.havingRaw(exp, "AND")
	return self
}

func (self *Builder) OrHavingRaw(exp string) *Builder {
	self.havingRaw(exp, "OR")
	return self
}

func (self *Builder) havingRaw(exp string, boolean string) {
	if len(self.components.Having) == 0 {
		boolean = ""
	}
	self.components.Having = append(self.components.Having, havingComponent{
		kind:    "raw",
		value:   exp,
		boolean: boolean,
	})
}

func (self *Builder) OrderBy(column string, direction string) *Builder {
	self.components.Order = append(self.components.Order, orderComponent{column, direction})
	return self
}

func (self *Builder) Limit(number interface{}) *Builder {
	self.components.Limit = []interface{}{number}
	self.bind("limit", number)
	return self
}

func (self *Builder) Offset(number interface{}) *Builder {
	self.components.Offset = []interface{}{number}
	self.bind("offset", number)
	return self
}

func (self *Builder) Delete() *Builder {
	if self.kind != "" {
		return self
	}
	self.kind = "delete"
	self.components.Delete = []interface{}{true}
	return self
}

func (self *Builder) Update(data Data) *Builder {
	if self.kind != "" {
		return self
	}

	self.kind = "update"
	self.components.Update = []interface{}{true}
	self.components.Set = []setComponent{setComponent(data)}

	self.bind("set", data.Values()...)

	return self
}

func (self *Builder) Insert(data ...Data) *Builder {
	if self.kind != "" && self.kind != "insert" {
		return self
	}

	self.kind = "insert"

	if len(data) == 0 {
		return self
	}

	columns := make([]interface{}, 0, len(data[0]))
	for _, c := range data[0].Keys() {
		columns = append(columns, c)
	}

	values := make([]valueComponent, 0, len(data))
	for _, d := range data {
		values = append(values, d.Values())
	}

	if len(self.components.Values) > 0 {
		self.components.Values = append(self.components.Values, values...)
	} else {
		self.components.Insert = []interface{}{true}
		self.components.Into = []interface{}{self.table}
		self.components.Columns = append(self.components.Columns, columns...)
		self.components.Values = append(self.components.Values, values...)
	}

	for _, v := range values {
		self.bind("values", v...)
	}

	return self
}

func (self *Builder) Count(column interface{}, alias string) *Builder {
	self.aggregate("count", column, alias)
	return self
}

func (self *Builder) Sum(column interface{}, alias string) *Builder {
	self.aggregate("sum", column, alias)
	return self
}

func (self *Builder) Avg(column interface{}, alias string) *Builder {
	self.aggregate("avg", column, alias)
	return self
}

func (self *Builder) Min(column interface{}, alias string) *Builder {
	self.aggregate("min", column, alias)
	return self
}

func (self *Builder) Max(column interface{}, alias string) *Builder {
	self.aggregate("max", column, alias)
	return self
}

func (self *Builder) aggregate(function string, column interface{}, alias string) {
	if self.kind != "" && self.kind != "select" {
		return
	}
	self.kind = "select"
	self.components.Aggregate = append(self.components.Aggregate, aggregateComponent{
		function: function,
		column:   column,
		alias:    alias,
	})
}

func (self *Builder) Sql() string {
	if driver == nil {
		panic("sqlx: driver is not defined")
	}
	if self.kind == "" {
		self.Select("*")
	}
	return newBaseGlammar(driver()).compile(self)
}

func (self *Builder) Data() []interface{} {
	if self.kind == "" {
		self.kind = "select"
	}
	bindings := make([]interface{}, 0)
	for _, k := range bindingsMap[self.kind] {
		for _, v := range self.bindings[k] {
			bindings = append(bindings, v)
		}
	}
	return bindings
}

func (self *Builder) bind(k string, b ...interface{}) {
	for _, v := range b {
		if _, ok := v.(Expression); !ok {
			self.bindings[k] = append(self.bindings[k], v)
		}
	}
}
