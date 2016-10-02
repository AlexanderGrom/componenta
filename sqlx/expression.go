package sqlx

type Expression struct {
	value    string
	bindings []interface{}
}

func (self Expression) String() string {
	return self.value
}

func (self Expression) Data() []interface{} {
	return self.bindings
}

func Raw(exp string, params ...interface{}) Expression {
	return Expression{
		value:    exp,
		bindings: params,
	}
}
