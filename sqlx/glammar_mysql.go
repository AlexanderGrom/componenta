package sqlx

type mysqlGlammar struct {
	baseGlammar
}

var _ glammar = (*mysqlGlammar)(nil)

func init() {
	registerDriver("mysql", newMysqlGlammar)
}

func newMysqlGlammar() glammar {
	g := &mysqlGlammar{}
	g.baseGlammar.glammar = g
	return g
}
