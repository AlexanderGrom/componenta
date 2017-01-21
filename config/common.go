package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Хранилище
var storage map[string]Config

// Ошибки
var (
	ErrPathEmpty      = errors.New("config: path empty")
	ErrPathNotCorrect = errors.New("config: path not correct")
	ErrConfigDontRead = errors.New("config: config file don`t read")
)

var regexpVarPath = regexp.MustCompile(`\$\{([A-Z_]+)\}`)

// Использование конкретного кофига
func Use(path string) (Config, error) {
	path, err := parsePath(path)
	if err != nil {
		return nil, err
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

// Перезагрузка конфиг файла
func Reload(path string) (Config, error) {
	path, err := parsePath(path)
	if err != nil {
		return nil, err
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

// Парсинг пути к файлу
func parsePath(path string) (string, error) {
	if len(path) == 0 {
		return "", ErrPathEmpty
	}
	path = regexpVarPath.ReplaceAllStringFunc(path, func(str string) string {
		return os.Getenv(str[2 : len(str)-1])
	})
	path, err := filepath.Abs(path)
	if err != nil {
		return "", ErrPathNotCorrect
	}
	return path, nil
}
