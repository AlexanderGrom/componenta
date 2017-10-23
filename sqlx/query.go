package sqlx

import (
	"database/sql"
)

type Querier interface {
	Query(*Builder) *Query
}

type DataBaser interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type DB struct {
	db *sql.DB
}

// Helper для добавления нового подключения
func DataBase(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

func (self *DB) Origin() *sql.DB {
	return self.db
}

func (self *DB) Query(builder *Builder) *Query {
	return &Query{
		db:    self.db,
		query: builder.Sql(),
		data:  builder.Data(),
	}
}

func (self *DB) QueryRaw(query string, data ...interface{}) *Query {
	return &Query{
		db:    self.db,
		query: query,
		data:  data,
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

func (self *Tx) Origin() *sql.Tx {
	return self.tx
}

func (self *Tx) Query(builder *Builder) *Query {
	return &Query{
		db:    self.tx,
		query: builder.Sql(),
		data:  builder.Data(),
	}
}

func (self *Tx) QueryRaw(query string, data ...interface{}) *Query {
	return &Query{
		db:    self.tx,
		query: query,
		data:  data,
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

type Result interface {
	LastInsertId() int64
	RowsAffected() int64
}

// Стандартный вариант Result предлагаемый драйвером к базе данных
type customResult struct {
	sql.Result
}

// Возвращает последний вставленный ID без проверки на поддержку драйвером
func (self customResult) LastInsertId() int64 {
	//id, _ := self.Result.LastInsertId()
	//return id
	return 0
}

// Возвращает кол-во затронутых строк последним запросом  без проверки на поддержку драйвером
func (self customResult) RowsAffected() int64 {
	//count, _ := self.Result.RowsAffected()
	//return count
	return 0
}

type Query struct {
	db    DataBaser
	query string
	data  []interface{}
}

// Выполнение запроса
func (self *Query) Exec() (Result, error) {
	res, err := self.db.Exec(self.query, self.data...)
	return customResult{res}, err
}

// Сканировать результаты
func (self *Query) Scan(a ...interface{}) error {
	rows, err := self.db.Query(self.query, self.data...)
	if err != nil {
		return err
	}
	return NewScanner(rows).Scan(a...)
}

// Сканировать в "чанки" и обрабатывать по кускам
func (self *Query) Chunk(i int, f ChunkFunk) error {
	rows, err := self.db.Query(self.query, self.data...)
	if err != nil {
		return err
	}
	return NewChunker(rows).Chunk(i, f)
}
