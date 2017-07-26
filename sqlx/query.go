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

func (self *Tx) Origin() *sql.Tx {
	return self.tx
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

type Result interface {
	LastInsertId() int64
	RowsAffected() int64
}

// Стандартный вариант Result предлагаемый драйвером к базе данных
type driverResult struct {
	sql.Result
}

// Возвращает последний вставленный ID без проверки на поддержку драйвером
func (self driverResult) LastInsertId() int64 {
	id, _ := self.Result.LastInsertId()
	return id
}

// Возвращает кол-во затронутых строк последним запросом  без проверки на поддержку драйвером
func (self driverResult) RowsAffected() int64 {
	count, _ := self.Result.RowsAffected()
	return count
}

// Пользовательский вариант Result, если драйвер к базы данных его не поддерживает
type customResult struct {
	lastId    int64
	rowsCount int64
}

// Возвращает последний вставленный ID
func (self customResult) LastInsertId() int64 {
	return self.lastId
}

// Возвращает кол-во затронутых строк последним запросом
func (self customResult) RowsAffected() int64 {
	return self.rowsCount
}

type Query struct {
	db      DataBaser
	builder *Builder
}

// Выполнение запроса
func (self *Query) Exec() (Result, error) {
	if self.builder.enableReturnId {
		rows, err := self.db.Query(self.builder.Sql(), self.builder.Data()...)
		if err != nil {
			return customResult{}, err
		}

		var id int64
		var count int64
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				return customResult{}, err
			}
			count++
		}

		if err := rows.Err(); err != nil {
			return customResult{}, err
		}
		if err := rows.Close(); err != nil {
			return customResult{}, err
		}

		return customResult{id, count}, err
	}
	res, err := self.db.Exec(self.builder.Sql(), self.builder.Data()...)
	return driverResult{res}, err
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
