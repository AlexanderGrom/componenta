package sqlx

type Expression struct {
	data string
}

func (self Expression) String() string {
	return self.data
}

func Raw(exp string) Expression {
	return Expression{
		data: exp,
	}
}
