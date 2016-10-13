package config

import (
	"bytes"
	"errors"
	"strings"
)

// Ошибки
var (
	ErrParseKey   = errors.New("сonfig: parse error (key)")
	ErrParseQuote = errors.New("сonfig: parse error (quote)")
)

// Парсер файла конфигурации
type parser struct {
	curPos  int           // Текущая позиция символа
	curChar rune          // Текущий символ (руна)
	textBuf []rune        // Буфер с рунами
	textLen int           // Длина буфера рун
	state   func()        // Текущее состояние автомата
	config  Config        // Собранный конфиг
	key     string        // Найденный ключ
	value   []interface{} // Найденные значения
}

// Конструктор
func newParser() *parser {
	return &parser{}
}

// Парсинг
func (self *parser) Parse(r []rune) Config {
	self.textBuf = r
	self.textLen = len(r)

	self.clean()
	self.movePos(0)
	self.saveState(self.matchKey)

	return self.extractСonfig()
}

// Очистка временных данных
func (self *parser) clean() {
	self.key = ""
	self.value = make([]interface{}, 0)
	self.config = make(Config)
}

// Установка нового состояния
func (self *parser) saveState(foo func()) {
	self.state = foo
}

// Получение следующего символа из входной строки
func (self *parser) moveNextPos() bool {
	return self.movePos(self.curPos + 1)
}

// Перемещает указатель на указанную позицию во входной строке и считывание символа
func (self *parser) movePos(pos int) bool {
	self.curPos = pos

	if pos < self.textLen && pos >= 0 {
		self.curChar = self.textBuf[pos]
	} else {
		self.curChar = 0
	}

	if self.curChar == 0 {
		return false
	}

	return true
}

// Определяет конец буфера с данными
func (self *parser) isEOF() bool {
	return self.curPos >= self.textLen
}

// Символы из которых может состоять ключ конфигурации
func (self *parser) isKey(char rune) bool {
	switch {
	case char >= '0' && char <= '9':
		return true
	case char >= 'A' && char <= 'Z':
		return true
	case char >= 'a' && char <= 'z':
		return true
	case char == '.' || char == '_':
		return true
	}
	return false
}

// Пробельные символы
func (self *parser) isSpace(char rune) bool {
	switch char {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}

// Символ конца строки (\n)
func (self *parser) isNL(char rune) bool {
	return char == '\n'
}

// Символ равенства
func (self *parser) isEq(char rune) bool {
	return char == '='
}

// Символы комментария
func (self *parser) isComment(char rune) bool {
	return char == '#'
}

// Символы кавычек
func (self *parser) isQuote(char rune) bool {
	switch char {
	case '"', '\'':
		return true
	}
	return false
}

// Символы разделения элементов в списке элементов
func (self *parser) isSeparator(char rune) bool {
	return char == ','
}

// Пропускает пробелы и символы перевода строк и возвращает кол-во пропущенных
func (self *parser) skipSpaces() int {
	count := 0
	for self.isSpace(self.curChar) {
		self.moveNextPos()
		count++
	}
	return count
}

// Пропускает комментарии до конца строки и возвращает кол-во пропущенных символов
func (self *parser) skipComment() int {
	count := 0
	if self.isComment(self.curChar) {
		count = self.skipLine()
	}
	return count
}

// Пропускает пробелы, переводы строк, табуляцию и комментарии. Всякий мусок вообщем...
func (self *parser) skipTrash() int {
	count := 0
	for {
		chars := 0
		chars += self.skipSpaces()
		chars += self.skipComment()
		count += chars
		if chars > 0 {
			continue
		}
		break
	}
	return count
}

// Пропускает текст до конца строки и возвращает кол-во пропущенных символов
func (self *parser) skipLine() int {
	count := 0
	for !self.isNL(self.curChar) && !self.isEOF() {
		self.moveNextPos()
		count++
	}
	return count
}

// Извлечение ключей и их значений из буфера данных
func (self *parser) extractСonfig() Config {
	for !self.isEOF() {
		self.state()
	}
	return self.config
}

// Поиск ключа
func (self *parser) matchKey() {
	self.skipTrash()

	buff := bytes.Buffer{}
	for self.isKey(self.curChar) {
		buff.WriteRune(self.curChar)
		self.moveNextPos()
	}

	self.key = buff.String()

	if len(self.key) == 0 && !self.isEOF() {
		panic(ErrParseKey)
	}

	self.saveState(self.matchEq)
}

// Поиск символа равенства, который следует за ключом
func (self *parser) matchEq() {
	self.skipSpaces()
	if self.isEq(self.curChar) {
		self.moveNextPos()
		self.saveState(self.matchValue)
	} else {
		self.saveState(self.matchKey)
	}
}

// Поиск значения
func (self *parser) matchValue() {
	self.skipSpaces()

	if self.isQuote(self.curChar) {
		quote := self.curChar
		escape := false

		self.moveNextPos()

		buff := bytes.Buffer{}
		for !self.isEOF() && (self.curChar != quote || escape == true) {
			buff.WriteRune(self.curChar)

			escape = (self.curChar == '\\')

			self.moveNextPos()
		}

		if self.curChar != quote {
			panic(ErrParseQuote)
		}

		self.value = append(self.value, buff.String())
		self.moveNextPos()
	} else {
		buff := bytes.Buffer{}
		for !self.isNL(self.curChar) && !self.isEOF() && !self.isSeparator(self.curChar) && !self.isComment(self.curChar) {
			buff.WriteRune(self.curChar)
			self.moveNextPos()
		}

		self.value = append(self.value, strings.TrimSpace(buff.String()))
	}

	if self.isSeparator(self.curChar) {
		self.saveState(self.matchValueList)
		self.moveNextPos()
	} else {
		self.saveState(self.saveConfig)
		self.skipLine()
	}
}

// Поиск значения если мы в состоянии списка значений
func (self *parser) matchValueList() {
	self.skipTrash()
	self.saveState(self.matchValue)
}

// Поиск значения если мы в состоянии списка значений
func (self *parser) saveConfig() {
	self.config[self.key] = self.value
	self.key = ""
	self.value = make([]interface{}, 0)
	self.saveState(self.matchKey)
}
