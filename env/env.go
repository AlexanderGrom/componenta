package env

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
)

// Ошибки
var (
	ErrPathEmpty    = errors.New("env: path empty")
	ErrFileDontRead = errors.New("env: file don`t read")
)

var regexpVarPath = regexp.MustCompile(`\$\{([A-Z_]+)\}`)

type Env map[string]string

// Загрузка файла в системные пременные
func Load(path string) error {
	path, err := parsePath(path)
	if err != nil {
		return err
	}
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return ErrFileDontRead
	}
	vars := NewParser().Parse(bytes.Runes(buffer))
	for key, value := range vars {
		os.Setenv(key, value)
	}
	return nil
}

// Парсинг пути к файлу
func parsePath(path string) (string, error) {
	if len(path) == 0 {
		return "", ErrPathEmpty
	}
	path = regexpVarPath.ReplaceAllStringFunc(path, func(str string) string {
		return os.Getenv(str[2 : len(str)-1])
	})
	return path, nil
}
