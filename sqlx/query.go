package sqlx

import (
	"database/sql"
)

type DataBaser interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type Transacter interface {
	Query(*Builder) *Query
	Commit() error
	Rollback() error
}

type DB struct {
	db *sql.DB
	dt bool
}

// Helper для добавления нового подключения
func DataBase(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

func (self *DB) Query(builder *Builder) *Query {
	return &Query{
		db:      self.db,
		builder: builder,
	}
}

// Включает/отключает механизм транзакций через Begin().
// Если toggle установлен в true, то вызовы Begin() будут возвращать обычный коннект к базе данных.
// В основном она нужно для тестов.
func (self *DB) DisableTransaction(toggle bool) {
	self.dt = toggle
}

// Начать транзакцию
func (self *DB) Begin() (*Tx, error) {
	if self.dt {
		return &Tx{&txdb{self.db}}, nil
	}
	tx, err := self.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{&txtx{tx}}, nil
}

type Tx struct {
	Transacter
}

// Реальные транзакции
type txtx struct {
	tx *sql.Tx
}

func (self *txtx) Query(builder *Builder) *Query {
	return &Query{
		db:      self.tx,
		builder: builder,
	}
}

// Зафиксировать транзакцию
func (self *txtx) Commit() error {
	return self.tx.Commit()
}

// Откатить транзакцию
func (self *txtx) Rollback() error {
	return self.tx.Rollback()
}

// Заглушка, если транзакции отключены через DisableTransaction()
type txdb struct {
	db *sql.DB
}

func (self *txdb) Query(builder *Builder) *Query {
	return &Query{
		db:      self.db,
		builder: builder,
	}
}

// Зафиксировать транзакцию
func (self *txdb) Commit() error {
	return nil
}

// Откатить транзакцию
func (self *txdb) Rollback() error {
	return nil
}

type Result struct {
	sql.Result
}

type Query struct {
	db      DataBaser
	builder *Builder
}

// Выполнение запроса
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
