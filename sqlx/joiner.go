package sqlx

type joinCondition struct {
	kind     string
	column   string
	operator string
	value    interface{}
	boolean  string
}

type Joiner struct {
	kind       string
	table      string
	conditions []joinCondition
	bindings   []interface{}
}

func newJoiner(table, kind string) *Joiner {
	return &Joiner{
		kind:  kind,
		table: table,
	}
}

func (self *Joiner) On(column1, operator, column2 string) *Joiner {
	self.where("on", column1, operator, column2, "AND")
	return self
}

func (self *Joiner) OrOn(column1, operator, column2 string) *Joiner {
	self.where("on", column1, operator, column2, "OR")
	return self
}

func (self *Joiner) Where(column string, operator string, value interface{}) *Joiner {
	self.where("where", column, operator, value, "AND")
	return self
}

func (self *Joiner) OrWhere(column string, operator string, value interface{}) *Joiner {
	self.where("where", column, operator, value, "OR")
	return self
}

func (self *Joiner) where(kind, column, operator string, value interface{}, boolean string) {
	if !isOperator(operator) {
		panic("sqlx: such a \"operator\" is not allowed")
	}
	if len(self.conditions) == 0 {
		boolean = ""
	}
	self.conditions = append(self.conditions, joinCondition{
		kind:     kind,
		column:   column,
		operator: operator,
		value:    value,
		boolean:  boolean,
	})
	if kind == "where" {
		self.bind(value)
	}
}

func (self *Joiner) bind(b ...interface{}) {
	for _, v := range b {
		if exp, ok := v.(Expression); ok {
			self.bindings = append(self.bindings, exp.Data()...)
		} else {
			self.bindings = append(self.bindings, v)
		}
	}
}
