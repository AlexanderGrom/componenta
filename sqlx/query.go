package sqlx

import (
	"database/sql"
)

type DataBaser interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type DB struct {
	db *sql.DB
}

// Helper для добавления нового подключения
func DataBase(db *sql.DB) *DB {
	return &DB{db}
}

func (self *DB) Query(builder *Builder) *Query {
	return &Query{
		db:      self.db,
		builder: builder,
	}
}

// Начать транзакцию
func (self *DB) Begin() (*Tx, error) {
	tx, err := self.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, err
}

type Tx struct {
	tx *sql.Tx
}

func (self *Tx) Query(builder *Builder) *Query {
	return &Query{
		db:      self.tx,
		builder: builder,
	}
}

// Зафиксировать транзакцию
func (self *Tx) Commit() error {
	return self.tx.Commit()
}

// Откатить транзакцию
func (self *Tx) Rollback() error {
	return self.tx.Rollback()
}

type Result struct {
	sql.Result
}

type Query struct {
	db      DataBaser
	builder *Builder
}

// Зафиксировать транзакцию
func (self *Query) Exec() (*Result, error) {
	res, err := self.db.Exec(self.builder.Sql(), self.builder.Data()...)
	if err != nil {
		return nil, err
	}
	return &Result{res}, err
}

// Сканировать результаты
func (self *Query) Scan(a ...interface{}) error {
	rows, err := self.db.Query(self.builder.Sql(), self.builder.Data()...)
	if err != nil {
		return err
	}
	return NewScanner(rows).Scan(a...)
}

// Сканировать в "чанки" и обрабатывать по кускам
func (self *Query) Chunk(i int, f ChunkFunk) error {
	rows, err := self.db.Query(self.builder.Sql(), self.builder.Data()...)
	if err != nil {
		return err
	}
	return NewChunker(rows).Chunk(i, f)
}
