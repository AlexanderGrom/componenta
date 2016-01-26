package i18n

import (
    "bytes"
)

//
// Парсер строки
//
type parser struct {
    curPos  int    // Текущая позиция символа
    curChar rune   // Текущий символ (руна)
    textBuf []rune // Буфер с рунами
    textLen int    // Длина буфера рун
}

//
// Конструктор
//
func newParser() *parser {
    return &parser{
        curPos:  -1,
        curChar: 0,
        textBuf: []rune{},
        textLen: 0,
    }
}

//
// Парсинг строки
// Вернет ключ, значение и логическое значение об успешности поиска
//
func (self *parser) Parse(text string) (key string, value interface{}, matched bool) {
    self.curPos = -1
    self.curChar = 0
    self.textBuf = []rune(text)
    self.textLen = len(self.textBuf)

    self.movePos(0)

    matched = self.matchConfig(&key, &value)

    return key, value, matched
}

//
// Получение следующего символа из входной строки
//
func (self *parser) moveNextPos() bool {
    return self.movePos(self.curPos + 1)
}

//
// Перемещает указатель на указанную позицию во входной строке и считывание символа
//
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

//
// Определяет конец строки
//
func (self *parser) isEOF() bool {
    return self.curPos >= self.textLen && self.textLen > 0
}

//
// Символы из которых может состоять ключ конфигурации
//
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

//
// Пробельные символы
//
func (self *parser) isSpace(char rune) bool {
    switch char {
    case ' ', '\t', '\r', '\n':
        return true
    }
    return false
}

//
// Символы комментария
//
func (self *parser) isComment(char rune) bool {
    return char == '#'
}

//
// Символы кавычек
//
func (self *parser) isQuote(char rune) bool {
    switch char {
    case '"', '\'':
        return true
    }
    return false
}

//
// Символы разделения элементов в списке элементов
//
func (self *parser) isSeparator(char rune) bool {
    switch char {
    case ' ', ',', '\r', '\n':
        return true
    }
    return false
}

//
// Пропускает пробелы и символы перевода строк и возвращает кол-во пропущенных
//
func (self *parser) skipSpaces() int {
    count := 0
    for self.isSpace(self.curChar) {
        self.moveNextPos()
        count++
    }
    return count
}

//
// Поиск конфигурации в строке
//
func (self *parser) matchConfig(key *string, value *interface{}) bool {
    self.skipSpaces()

    if self.isEOF() {
        return false
    }

    if self.isComment(self.curChar) {
        return false
    }

    if !self.matchKey(key) {
        return false
    }

    self.skipSpaces()

    if self.curChar == '=' {
        self.moveNextPos()
    } else {
        return true
    }

    self.skipSpaces()

    val := ""
    arr := make([]interface{}, 0)
    for self.matchValue(&val) {
        arr = append(arr, val)
        val = ""
    }

    switch {
    case len(arr) == 1:
        *value = arr[0]
    case len(arr) > 1:
        *value = arr
    }

    return true
}

//
// Поиск ключа
//
func (self *parser) matchKey(key *string) bool {
    buff := bytes.Buffer{}
    for self.isKey(self.curChar) {
        buff.WriteRune(self.curChar)
        self.moveNextPos()
    }
    *key = buff.String()

    if len(*key) == 0 {
        return false
    }

    return true
}

//
// Поиск значения
// Значение может быть в кавычках
//
func (self *parser) matchValue(value *string) bool {
    if self.isSpace(self.curChar) {
        self.skipSpaces()
    }

    if self.isComment(self.curChar) {
        return false
    }

    if self.isEOF() {
        return false
    }

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

        *value = buff.String()

        if self.curChar != quote {
            return false
        }

        self.moveNextPos()
    } else {
        buff := bytes.Buffer{}
        for !self.isEOF() && !self.isSeparator(self.curChar) && !self.isComment(self.curChar) {
            buff.WriteRune(self.curChar)
            self.moveNextPos()
        }

        *value = buff.String()
    }

    self.skipSpaces()

    if self.isSeparator(self.curChar) {
        self.moveNextPos()
    }

    return true
}
