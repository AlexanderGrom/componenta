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

// Вставка Insert IGNORE
func (self *mysqlGlammar) compileOrIgnore(b *Builder) string {
	if len(b.components.OrIgnore) == 0 {
		return ""
	}
	return "IGNORE"
}
