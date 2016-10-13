package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
)

// Хранилище
var storage map[string]Config

// Ошибки
var (
	ErrPathEmpty      = errors.New("сonfig: path empty")
	ErrPathNotCorrect = errors.New("сonfig: path not correct")
	ErrConfigDontRead = errors.New("сonfig: config file don`t read")
)

// Использование конкретного кофига
func Use(path string) (Config, error) {
	if len(path) == 0 {
		return nil, ErrPathEmpty
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, ErrPathNotCorrect
	}
	if cfg, ok := storage[path]; ok {
		return cfg, nil
	}
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrConfigDontRead
	}
	return parse(bytes.Runes(buffer)), nil
}

// Парсинг
func parse(r []rune) Config {
	return newParser().Parse(r)
}
