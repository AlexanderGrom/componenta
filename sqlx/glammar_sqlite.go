package sqlx

type sqliteGlammar struct {
	baseGlammar
}

var _ glammar = (*sqliteGlammar)(nil)

func init() {
	registerDriver("sqlite3", newSqliteGlammar)
}

func newSqliteGlammar() glammar {
	g := &sqliteGlammar{}
	g.baseGlammar.glammar = g
	return g
}

// Вставка Insert OR IGNORE
func (self *sqliteGlammar) compileOrIgnore(b *Builder) string {
	if len(b.components.OrIgnore) == 0 {
		return ""
	}
	return "OR IGNORE"
}
